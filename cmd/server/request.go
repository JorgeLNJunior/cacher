package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	Operation Operation
	Key       string
	Value     string
	Expiry    time.Time
}

type Operation string

func (o Operation) Valid() bool {
	switch {
	case o == OperationGet:
		return true
	case o == OperationSet:
		return true
	case o == OperationDel:
		return true
	case o == OperationExp:
		return true
	default:
		return false
	}
}

func (o Operation) String() string {
	return string(o)
}

const (
	OperationGet Operation = "GET"
	OperationSet Operation = "SET"
	OperationDel Operation = "DEL"
	OperationExp Operation = "EXP"
)

const maxParameters = 3

var (
	ErrInvalidOperation     = errors.New("operation must be GET, SET or DEL")
	ErrInvalidFormat        = errors.New("message format does not complain")
	ErrNoKey                = errors.New("should provide a key")
	ErrNoValue              = errors.New("should provide a value when operation is SET")
	ErrInvalidUnixTimestamp = errors.New("should provide a valid a unix timestamp")
)

func (r *Request) Marshal() ([]byte, error) {
	if !r.Operation.Valid() {
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
	trimData := strings.TrimSuffix(string(data), "\n") // messages are ending with a \n and we should remove it
	splitData := strings.SplitN(trimData, " ", maxParameters)
	if len(splitData) < 2 {
		return ErrInvalidFormat
	}

	operation := Operation(splitData[0])
	if !operation.Valid() {
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

	if operation == OperationGet || operation == OperationDel {
		r.Operation = operation
		r.Key = splitData[1]
		return nil
	}

	if operation == OperationExp {
		if len(splitData) < 3 {
			return ErrInvalidFormat
		}

		r.Operation = operation
		r.Key = splitData[1]

		seconds, err := strconv.ParseInt(splitData[2], 10, 64)
		if err != nil {
			return ErrInvalidUnixTimestamp
		}

		exp := time.Unix(seconds, 0)
		if time.Now().After(exp) {
			return ErrInvalidUnixTimestamp
		}
		r.Expiry = exp

		return nil
	}

	return errors.New("unexpected error")
}

func (r Request) String() string {
	v := string(r.Operation) + " " + r.Key
	if len(r.Value) > 0 {
		v += " " + r.Value
	}
	return v
}
