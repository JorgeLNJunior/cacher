package main

import (
	"errors"
)

type Request struct {
	Operation Operation
	Key       string
	Value     []byte
}

type Operation string

const (
	OperationGet Operation = "GET"
	OperationSet Operation = "SET"
)

var (
	ErrInvalidOperation = errors.New("operation must be GET or SET")
	ErrNoKey            = errors.New("should provide a key when operation is GET")
	ErrNoValue          = errors.New("should provide a value when operation is SET")
)

func (m *Request) Marshal() ([]byte, error) {
	data := make([]byte, 0)

	if m.Operation == "" {
		return nil, ErrInvalidOperation
	}
	if m.Operation != OperationGet && m.Operation != OperationSet {
		return nil, ErrInvalidOperation
	}
	if m.Key == "" {
		return nil, ErrNoKey
	}
	if m.Operation == OperationSet && len(m.Value) < 1 {
		return nil, ErrNoValue
	}

	writeToByteSlice([]byte(m.Operation), data)

	for _, v := range "Key: " {
		data = append(data, byte(v))
	}
	writeToByteSlice([]byte(m.Key), data)

	if m.Operation == "SET" {
		for _, v := range "Value: " {
			data = append(data, byte(v))
		}
		writeToByteSlice(m.Value, data)
	}

	return data, nil
}

func (m *Request) Unmarshal(data []byte) error {
	return nil
}

// writeToByteSlice copy the values from a slice of bytes to
// another slice of bytes, appending a '\n' to the seccond slice.
func writeToByteSlice(from []byte, to []byte) {
	for i := 0; i < len(from); i++ {
		to = append(to, byte(from[i]))

		isLastByte := i == (len(from) - 1)
		if isLastByte {
			to = append(to, byte('\n'))
		}
	}
}
