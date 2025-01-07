package main

import (
	"errors"
	"testing"
)

func TestResponseMarshal(t *testing.T) {
	t.Run("should return InvalidResponseStatus if status is invalid", func(tt *testing.T) {
		res := Response{
			Status: "INVALID",
		}

		if _, err := res.Marshal(); !errors.Is(err, ErrInvalidResponseStatus) {
			tt.Errorf("expected ErrInvalidResponseStatus but received '%s'", err)
		}
	})

	t.Run("should marshal the data", func(tt *testing.T) {
		res := Response{
			Status:  ResponseStatusOK,
			Message: "OK",
		}

		if _, err := res.Marshal(); err != nil {
			tt.Error(err)
		}
	})
}

func TestResponseUnmarshal(t *testing.T) {
	t.Run("should return InvalidResponseStatus if status is smaller than 2 bytes", func(tt *testing.T) {
		res := Response{
			Status: "I",
		}

		if _, err := res.Marshal(); !errors.Is(err, ErrInvalidResponseStatus) {
			tt.Errorf("expected ErrInvalidResponseStatus but received '%s'", err)
		}
	})

	t.Run("should return InvalidResponseStatus if status is invalid", func(tt *testing.T) {
		res := Response{
			Status: "INVALID",
		}

		if _, err := res.Marshal(); !errors.Is(err, ErrInvalidResponseStatus) {
			tt.Errorf("expected ErrInvalidResponseStatus but received '%s'", err)
		}
	})

	t.Run("should unmarshal the data", func(tt *testing.T) {
		data := []byte("OKthis is a message")

		res := Response{}
		if err := res.Unmarshal(data); err != nil {
			tt.Error(err)
		}
	})
}
