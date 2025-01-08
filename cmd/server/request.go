package main

import (
	"errors"
	"strings"
)

type Request struct {
	Operation Operation
	Key       string
	Value     string
}

type Operation string

func (o Operation) String() string {
	return string(o)
}

const (
	OperationGet Operation = "GET"
	OperationSet Operation = "SET"
)

const maxParameters = 3

var (
	ErrInvalidOperation = errors.New("operation must be GET or SET")
	ErrInvalidFormat    = errors.New("message format does not complain")
	ErrNoKey            = errors.New("should provide a key when operation is GET")
	ErrNoValue          = errors.New("should provide a value when operation is SET")
)

func (r *Request) Marshal() ([]byte, error) {
	if r.Operation != OperationGet && r.Operation != OperationSet {
		return nil, ErrInvalidOperation
	}
	if r.Key == "" {
		return nil, ErrNoKey
	}
	if r.Operation == OperationSet && len(r.Value) < 1 {
		return nil, ErrNoValue
	}

	data := r.Operation.String() + " " + r.Key
	if len(r.Value) > 0 {
		data += " " + r.Value
	}

	return []byte(data), nil
}

func (r *Request) Unmarshal(data []byte) error {
	splitData := strings.SplitN(string(data), " ", maxParameters)
	if len(splitData) < 2 {
		return ErrInvalidFormat
	}

	operation := Operation(splitData[0])
	if operation != OperationGet && operation != OperationSet {
		return ErrInvalidOperation
	}

	if operation == OperationSet {
		if len(splitData) < 3 {
			return ErrInvalidFormat
		}

		r.Operation = operation
		r.Key = splitData[1]
		r.Value = splitData[2]
		return nil
	}

	if operation == OperationGet {
		r.Operation = operation
		r.Key = splitData[1]
		return nil
	}

	return errors.New("unexpected error")
}

func (r Request) String() string {
	v := r.Operation.String() + " " + r.Key
	if len(r.Value) > 0 {
		v += " " + r.Value
	}
	return v
}
