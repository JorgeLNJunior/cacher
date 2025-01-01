package main

import (
	"errors"
	"testing"
)

func TestMarshal(t *testing.T) {
	t.Run("should return an error if Operation is empty", func(tt *testing.T) {
		req := Request{
			Operation: "",
		}

		if _, err := req.Marshal(); !errors.Is(err, ErrInvalidOperation) {
			tt.Fatalf("expected 'ErrInvalidOperation' but received '%s'", err)
		}
	})

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
}
