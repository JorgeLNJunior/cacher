package main

import (
	"errors"
	"strings"
	"testing"
)

func TestMarshal(t *testing.T) {
	t.Run("should return an error if Operation has an invalid value", func(tt *testing.T) {
		req := Request{
			Operation: "PUT",
		}

		if _, err := req.Marshal(); !errors.Is(err, ErrInvalidOperation) {
			tt.Fatalf("expected 'ErrInvalidOperation' but received '%s'", err)
		}
	})

	t.Run("should return an error if Key is empty", func(tt *testing.T) {
		req := Request{
			Operation: OperationGet,
			Key:       "",
		}

		if _, err := req.Marshal(); !errors.Is(err, ErrNoKey) {
			tt.Fatalf("expected 'ErrNoKey' but received '%s'", err)
		}
	})

	t.Run("should return an error if Operation is SET but the Value is empty", func(tt *testing.T) {
		req := Request{
			Operation: OperationSet,
			Key:       "foo",
		}

		if _, err := req.Marshal(); !errors.Is(err, ErrNoValue) {
			tt.Fatalf("expected 'ErrNoValue' but received '%s'", err)
		}
	})

	t.Run("should marshal a GET operation", func(tt *testing.T) {
		req := Request{
			Operation: OperationGet,
			Key:       "foo",
		}

		data, err := req.Marshal()
		if err != nil {
			tt.Fatal(err)
		}

		dataSplit := strings.Split(string(data), " ")
		if len(dataSplit) < 2 {
			tt.Errorf("expected data to have a operation and a key but got '%s'", string(data))
		}

		operation := dataSplit[0]
		key := dataSplit[1]
		if operation != req.Operation.String() {
			tt.Errorf("expected operation to be '%s' but got '%s'", req.Operation.String(), operation)
		}
		if key != req.Key {
			tt.Errorf("expected key to be '%s' but got '%s'", req.Key, key)
		}
	})
}

func TestUnmarshal(t *testing.T) {
	t.Run("should return an error if the format is invalid", func(tt *testing.T) {
		bytes := []byte("SET ") // expected Key to have a value

		result := Request{}
		if err := result.Unmarshal(bytes); !errors.Is(err, ErrInvalidFormat) {
			tt.Errorf("expected ErrInvalidFormat but received %s", err)
		}
	})

	t.Run("should return an error if the operation is invalid", func(tt *testing.T) {
		bytes := []byte("INVALID key")

		result := Request{}
		if err := result.Unmarshal(bytes); !errors.Is(err, ErrInvalidOperation) {
			tt.Errorf("expected ErrInvalidOperation but received %s", err)
		}
	})

	t.Run("should return an error if Operation is SET and Value was not provided", func(tt *testing.T) {
		bytes := []byte("SET key")

		result := Request{}
		if err := result.Unmarshal(bytes); !errors.Is(err, ErrInvalidFormat) {
			tt.Errorf("expected ErrInvalidFormat but received %s", err)
		}
	})

	t.Run("should marshal a GET operation", func(tt *testing.T) {
		operation := OperationGet
		key := "foo"
		bytes := []byte(string(operation) + " " + key)

		result := Request{}
		if err := result.Unmarshal(bytes); err != nil {
			tt.Error(err)
		}

		if result.Operation != operation {
			tt.Errorf("expected operation to be '%s' but got '%s'", operation, result.Operation)
		}
		if result.Key != key {
			tt.Errorf("expected key to be '%s' but got '%s'", key, result.Key)
		}
		if len(result.Value) > 0 {
			tt.Errorf("expected a 0 length value but got a length of %d", len(result.Value))
		}
	})

	t.Run("should marshal a SET operation", func(tt *testing.T) {
		operation := OperationSet
		key := "foo"
		value := "bar"
		bytes := []byte(string(operation) + " " + key + " " + value)

		result := Request{}
		if err := result.Unmarshal(bytes); err != nil {
			tt.Error(err)
		}

		if result.Operation != operation {
			tt.Errorf("expected operation to be '%s' but got '%s'", operation, result.Operation)
		}
		if result.Key != key {
			tt.Errorf("expected key to be '%s' but got '%s'", key, result.Key)
		}
		if result.Value != value {
			tt.Errorf("expected value to be '%s' but got '%s'", value, result.Value)
		}
	})
}
