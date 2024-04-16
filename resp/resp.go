package resp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)



type RESPType byte

// Constants of type RESPType representing the RESP protocol prefixes.
const (
    SimpleStringPrefix RESPType = '+'
    ErrorPrefix        RESPType = '-'
    IntegerPrefix      RESPType = ':'
    BulkStringPrefix   RESPType = '$'
    ArrayPrefix        RESPType = '*'
)

type Payload struct {
	SimpleString string
	Error string
	Integer int
	BulkString string
	Array []interface{}
	Type RESPType
}

type Resp struct{
	str string
	payload Payload
}

func readerCurrentValue(reader *bufio.Reader){
	currValue, err := reader.Peek(1)
	if err!=nil {
		fmt.Println("error while peaking",err)
	}
	fmt.Println("current reader value",string(rune(currValue[0])))
}

func (r *Resp) parseString(reader *bufio.Reader) (Payload, error) {
	readerCurrentValue(reader)
	typeByte, err := reader.ReadByte()
	if err != nil {
		return Payload{}, err
	}

	typePayload := RESPType(typeByte)
	

	switch typePayload {
	case SimpleStringPrefix:
		return r.parseSimpleString(reader)
	case ErrorPrefix:
		return r.parseErrorString(reader)
	case IntegerPrefix:
		return r.parseIntegerString(reader)
	case BulkStringPrefix:
		return r.parseBulkString(reader)
	case ArrayPrefix:
		return r.parseArrayString(reader)
	default:
		return Payload{}, fmt.Errorf("unknown RESP type prefix: %v", typeByte)
	}
}

func readTrimmedString(reader *bufio.Reader) (string, error) {
    line, err := reader.ReadString('\n')
	
    if err != nil {
        return "", err
    }
	line = strings.TrimSuffix(line, "\r\n")
	
	return line,nil
}



func (r *Resp) parseSimpleString(reader *bufio.Reader) (Payload, error) {
	line, err := readTrimmedString(reader)
	if err!= nil{
		return Payload{}, err
	}
	r.payload.SimpleString = line
	r.payload.Type = SimpleStringPrefix

	return r.payload, nil

}

func (r *Resp) parseErrorString(reader *bufio.Reader) (Payload, error) {
	line, err := readTrimmedString(reader)
	if err!= nil{
		return Payload{}, err
	}
	r.payload.Error = line
	r.payload.Type = ErrorPrefix

	return r.payload, nil
}

func (r *Resp) parseIntegerString(reader *bufio.Reader) (Payload, error){
	line, err := readTrimmedString(reader)
	if err!= nil{
		return Payload{}, err
	}
	integer, err := strconv.Atoi(line)
	if err != nil {
		return Payload{}, err
	}
	r.payload.Integer = integer
	r.payload.Type = IntegerPrefix
	return r.payload, nil
}

func (r *Resp) parseBulkString(reader *bufio.Reader) (Payload,error) {
	line, err := readTrimmedString(reader)
	if err!= nil{
		return Payload{}, err
	}
	length, err := strconv.Atoi(line)
	if err != nil {
		return Payload{}, err
	}

	
	if length == -1 {
		r.payload.BulkString = ""
		r.payload.Type = BulkStringPrefix
		reader.ReadString('\n')
		return r.payload, nil
	}
	
	bulkStringBytes := make([]byte, length)
	_, err = io.ReadFull(reader, bulkStringBytes)
	if err != nil {
		return Payload{}, err
	}
	r.payload.BulkString = string(bulkStringBytes)
	r.payload.Type = BulkStringPrefix
	reader.ReadString('\n')

	return r.payload, nil
}

func(r *Resp) parseArrayString(reader *bufio.Reader)(Payload, error){
	line, err := readTrimmedString(reader)
	if err!= nil{
		return Payload{}, err
	}
	
	lengthArray, err := strconv.Atoi(line)
	
	if err != nil {
		return Payload{}, err
	}
	if lengthArray == -1 {
		r.payload.Array = nil
		r.payload.Type = ArrayPrefix
		return r.payload, nil
	}
	var array []interface{}
	for i:=0;i<lengthArray;i++ {
		element, err := r.parseString(reader)
		if err!= nil {
			return Payload{}, err
		}
		array = append(array,element)
		
	}
	r.payload.Array = array
	r.payload.Type = ArrayPrefix
	return r.payload, nil
}

func (r *Resp) SerializeResp() {

}

func (r *Resp) DeserializeResp(con net.Conn) (Payload, error){
	reader := bufio.NewReader(con)
	return r.parseString(reader)
}
