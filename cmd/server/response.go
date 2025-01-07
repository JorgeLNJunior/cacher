package main

import "errors"

type ResponseStatus string

const (
	ResponseStatusOK    = "OK"
	ResponseStatusError = "ER"
)

var ErrInvalidResponseStatus = errors.New("status must be OK or ER")

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

	data = append(data, byte(r.Status[0]), byte(r.Status[1]))
	for _, b := range r.Message {
		data = append(data, byte(b))
	}

	return data, nil
}

func (r *Response) Unmarshal(data []byte) error {
	if len(data) < 2 {
		return ErrInvalidResponseStatus
	}

	status := ResponseStatus(data[:2])
	if status != ResponseStatusOK && status != ResponseStatusError {
		return ErrInvalidResponseStatus
	}
	r.Status = status

	if len(data) > 2 {
		r.Message = string(data[:3])
	}

	return nil
}
