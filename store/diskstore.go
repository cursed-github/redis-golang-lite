package store

import (
	"bufio"
	"bytes"
	"fmt"
	"kvstore/resp"
	"os"
	"strconv"
	"strings"
	"time"
)

func WriteToDisk(payload resp.Payload)  error {
	data ,err:= SerializeData()
	if err!=nil {
		return err
	}
	err = saveToFile(data)
	if err!=nil {
		return err
	}

	return nil
}

func SerializeData() ([]byte, error){
	mutex.RLock()
	defer mutex.RUnlock()

	var buffer bytes.Buffer

	for key, val := range stringMap {
		if !val.ttl.IsZero() && time.Now().After(val.ttl) {
			continue
		}
		_, err := buffer.WriteString(fmt.Sprintf("%v,%v,%v\n", key, val.value, int64(time.Until(val.ttl).Seconds())))
		if err!=nil {
			return nil, err
		}
	}

	return buffer.Bytes(),nil
}

func saveToFile(data []byte) error {
    tempFilePath := "datastore.tmp"
    finalFilePath := "datastore.rdb"

    // Write data to a temporary file
    err := os.WriteFile(tempFilePath, data, 0644)
    if err != nil {
        return err
    }

    // Rename temporary file to final file name atomically
    err = os.Rename(tempFilePath, finalFilePath)
    if err != nil {
        return err
    }
    return nil
}

func ReadFromDisk() error {
	finalFilePath := "database.rdb"

	file, err := os.Open(finalFilePath)
	if err!= nil {
		return err
	}
	mutex.Lock() // Lock the mutex since we are modifying the global map
    defer mutex.Unlock()

	scanner := bufio.NewScanner(file)
	stringMap = make(map[string]stringvalue)

	for scanner.Scan() {
		line:= scanner.Text()
		lineArray := strings.Split(line, ",")
		if len(lineArray) != 3 {
            continue // handle error or malformed line
        }
		key:= lineArray[0]
		value:= lineArray[1]

		ttlSeconds, err := strconv.ParseInt(lineArray[2], 10, 64)
        if err != nil {
            return err // handle conversion error
        }

		var ttl time.Time
		if ttlSeconds>0 {
			ttl= time.Now().Add(time.Duration(ttlSeconds)*time.Second)
		}
		stringMap[key] = stringvalue{
			value: value,
			ttl: ttl,
		}

	
	}
	if err := scanner.Err(); err != nil {
        return err
    }

	return nil

}