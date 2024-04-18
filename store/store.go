package store

import (
	"fmt"
	"kvstore/resp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var stringMap = map[string]stringvalue{}
var mutex sync.RWMutex

type stringvalue struct{
	value string
	ttl time.Time
}

func ProcessRequest(payload resp.Payload) string {
	responsePayload := processPayload(payload)
	fmt.Println("payload processed got this response",responsePayload)
	responseString := resp.SerializeResp(responsePayload)
	fmt.Printf("Debug Output: %q\n", responseString)
	return responseString
}

func processPayload(payload resp.Payload) (resp.Payload) {
	if payload.Type != resp.ArrayPrefix {
        return errorResponse("Incorrect Message Type, Expected Array")
    }
	command:=payload.Array[0]

	switch  strings.ToLower(command.BulkString){
	case "ping":
		return ping(payload)
	case "echo":
		return echo(payload)
	case "set":
		return set(payload)
	case "get":
		return get(payload)	
	case "exists":
		return exist(payload)
	case "del":
		return del(payload)	
	case "save":
		return save(payload)
	case "config":
		return SimpleStringResponse("OK CONGIG")		
	}


	return errorResponse("Unsupported command")
}

func ping(payload resp.Payload) resp.Payload {
	return resp.Payload{
		BulkString: "PONG",
		Type: resp.BulkStringPrefix,
	}
}

func echo(payload resp.Payload) resp.Payload {
	return resp.Payload{
		BulkString: payload.Array[1].BulkString,
		Type: resp.BulkStringPrefix,
	}	
}

func set(payload resp.Payload) resp.Payload {
	if len(payload.Array) < 3 {
        return errorResponse("Incorrect number of arguments for 'set'")
    }

	key:= payload.Array[1].BulkString
	value:= payload.Array[2].BulkString

	if key=="key:__rand_int__" || value=="" {
		return resp.Payload{
			Error: "Incorrect Message Type, key or value empty",
			Type: resp.ErrorPrefix,
		}
	}
	fmt.Println("key value for redis-call",key, value)
	var ttl *time.Time
	if len(payload.Array)==5 {
		var err error
		ttl, err = getTTL(payload.Array[3].BulkString,payload.Array[4].BulkString)
		if err!=nil {
			return errorResponse(err.Error())
		}
	}

	mutex.Lock()
	defer mutex.Unlock()

	stringValue := stringvalue{
        value: value,
    }
    if ttl != nil {
        stringValue.ttl = *ttl
    }

    stringMap[key] = stringValue

	return SimpleStringResponse("OK")
}

func getTTL(ex_com, ex_value string) (*time.Time, error) {
	if ex_com == "" || ex_value =="" {
		return nil, fmt.Errorf("error with expirty command or expiry value")	
	}

	duration, err:= strconv.Atoi(ex_value)
	if err != nil {
		return nil, fmt.Errorf("unkown time expiry value%v",ex_com)	
	}
	var multiplier time.Duration
	var t time.Time
	switch ex_com {
	case "EX":
		multiplier = time.Second
		t= time.Now().Add(time.Duration(duration)*multiplier)
	case "PX":
		multiplier = time.Millisecond	
		t= time.Now().Add(time.Duration(duration)*multiplier)
	case "EXAT":
		t=time.Unix(int64(duration),0)	
	case "PXAT":
		t=time.Unix(0,int64(duration)*int64(time.Millisecond))	
	default:
		return nil, fmt.Errorf("unkown time expiry command%v",ex_com)	
	}
	
	return &t,nil
}

func get(payload resp.Payload) resp.Payload{
	if len(payload.Array) < 2 || payload.Array[1].BulkString == "" {
        return errorResponse("Problem with array or key is empty")
    }
	fmt.Println("get erquest recived")
	mutex.RLock()
	defer mutex.RUnlock()

	key:=  payload.Array[1].BulkString

	if stringValue, ok := stringMap[key] ; ok {
		if !stringValue.ttl.IsZero() && time.Now().After(stringValue.ttl){
			mutex.Lock()
			delete(stringMap,key)
			mutex.Unlock()
			return errorResponse("Key has expired and was deleted")
		}
		return SimpleStringResponse(stringValue.value)
	} 
	return errorResponse("Key does not exist")
}


func exist(payload resp.Payload) resp.Payload{
	if len(payload.Array)<2 || payload.Array[1].BulkString==""{
		return errorResponse("Error not enough keys to check for exist in input array")
	}
	var count int
	mutex.RLock()
	defer mutex.RUnlock()

	for i:=1;i<len(payload.Array);i++{
		key:= payload.Array[i].BulkString

		if stringvalue,ok := stringMap[key]; ok {
			if !stringvalue.ttl.IsZero() && time.Now().After(stringvalue.ttl){
				mutex.Lock()
				delete(stringMap,key)
				mutex.Unlock()
			} else {
				count++
			}
		}
	}
	return IntegerResponse(count)
}

func del(payload resp.Payload) resp.Payload{
	if len(payload.Array)<2 || payload.Array[1].BulkString==""{
		return errorResponse("Error not enough keys to check for exist in input array")
	}
	var count int
	mutex.Lock()
	defer mutex.Unlock()

	for i:=1;i<len(payload.Array);i++{
		key:= payload.Array[i].BulkString

		if _,ok := stringMap[key]; ok {
			delete(stringMap,key)
			count++
		}
	}
	return IntegerResponse(count)
}

func save(payload resp.Payload) resp.Payload {
	err := WriteToDisk(payload)
	if err!=nil {
		return errorResponse("failure while saving data to disk, check logs")
	}
	return SimpleStringResponse("OK")
}

func errorResponse(message string) resp.Payload {
    return resp.Payload{
        Error: message,
        Type: resp.ErrorPrefix,
    }
}

func SimpleStringResponse(message string) resp.Payload {
	return resp.Payload{
		SimpleString: message,
		Type: resp.SimpleStringPrefix,
	}
}

func IntegerResponse(message int) resp.Payload {
	return resp.Payload{
		Integer: message,
		Type: resp.IntegerPrefix,
	}
}