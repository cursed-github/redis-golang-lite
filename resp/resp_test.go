package resp

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseSimpleString(t *testing.T) {
	input := "+OK\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	resp := Resp{}

	payload, err := resp.parseString(reader)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if payload.SimpleString != "OK" || payload.Type != SimpleStringPrefix {
		t.Errorf("Expected 'OK' with type '+', got '%s' with type '%v'", payload.SimpleString, payload.Type)
	}
}

func TestParseErrorString(t *testing.T) {
	input := "-Error message\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	resp := Resp{}

	payload, err := resp.parseString(reader)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if payload.Error != "Error message" || payload.Type != ErrorPrefix {
		t.Errorf("Expected 'Error message' with type '-', got '%s' with type '%v'", payload.Error, payload.Type)
	}
}

func TestParseInteger(t *testing.T) {
	input := ":1000\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	resp := Resp{}

	payload, err := resp.parseString(reader)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if payload.Integer != 1000 || payload.Type != IntegerPrefix {
		t.Errorf("Expected 1000 with type ':', got %d with type '%v'", payload.Integer, payload.Type)
	}
}

func TestParseBulkString(t *testing.T) {
	input := "$6\r\nfoobar\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	resp := Resp{}

	payload, err := resp.parseString(reader)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if payload.BulkString != "foobar" || payload.Type != BulkStringPrefix {
		t.Errorf("Expected 'foobar' with type '$', got '%s' with type '%v'", payload.BulkString, payload.Type)
	}
}

func TestParseArray(t *testing.T) {
	input := "*2\r\n+OK\r\n:1000\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	resp := Resp{}

	payload, err := resp.parseString(reader)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(payload.Array) != 2 || payload.Type != ArrayPrefix {
		t.Fatalf("Expected array length 2 with type '*', got length %d with type '%v'", len(payload.Array), payload.Type)
	}
	// Asserting elements inside the array
	firstElement, ok := payload.Array[0].(Payload)
	if !ok || firstElement.SimpleString != "OK" || firstElement.Type != SimpleStringPrefix {
		t.Errorf("Expected first element 'OK' with type '+', got '%v' with type '%v'", firstElement.SimpleString, firstElement.Type)
	}
	secondElement, ok := payload.Array[1].(Payload)
	if !ok || secondElement.Integer != 1000 || secondElement.Type != IntegerPrefix {
		t.Errorf("Expected second element 1000 with type ':', got %v with type '%v'", secondElement.Integer, secondElement.Type)
	}
}

