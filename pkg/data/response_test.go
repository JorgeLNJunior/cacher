package data

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

	t.Run("should marshal the response", func(tt *testing.T) {
		res := Response{
			Status:  ResponseStatusOK,
			Message: "OK",
		}

		data, err := res.Marshal()
		if err != nil {
			tt.Error(err)
		}

		if string(data) != res.String() {
			tt.Errorf("expect a valid response but got '%s'", data)
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
