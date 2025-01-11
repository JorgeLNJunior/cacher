package data

import (
	"errors"
	"strings"
)

type ResponseStatus string

func (s ResponseStatus) String() string {
	return string(s)
}

const (
	ResponseStatusOK    = "OK"
	ResponseStatusError = "ERROR"
)

var ErrInvalidResponseStatus = errors.New("status must be OK or ERROR")

type Response struct {
	Status  ResponseStatus
	Message string
}

func NewResponse(status ResponseStatus, msg string) Response {
	return Response{status, msg}
}

func (r Response) Marshal() ([]byte, error) {
	data := make([]byte, 0)

	if r.Status != ResponseStatusOK && r.Status != ResponseStatusError {
		return nil, ErrInvalidResponseStatus
	}

	for _, v := range r.Status {
		data = append(data, byte(v))
	}
	data = append(data, byte(' '))
	for _, b := range r.Message {
		data = append(data, byte(b))
	}

	return data, nil
}

func (r *Response) Unmarshal(data []byte) error {
	splitData := strings.SplitN(string(data), " ", 2)
	if len(splitData) < 2 {
		return ErrInvalidResponseStatus
	}

	status := ResponseStatus(splitData[0])
	if status != ResponseStatusOK && status != ResponseStatusError {
		return ErrInvalidResponseStatus
	}
	r.Status = status
	r.Message = string(splitData[1])

	return nil
}

func (r Response) String() string {
	return r.Status.String() + " " + r.Message
}
