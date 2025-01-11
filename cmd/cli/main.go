package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/JorgeLNJunior/cacher/pkg/data"
)

func main() {
	var operation string
	var key string
	var value string
	var expiry int64
	var url string

	flag.StringVar(&operation, "operation", data.OperationGet.String(), "the operation to be done. GET, SET, DEL or EXP")
	flag.StringVar(&key, "key", "", "the key to send in the request")
	flag.StringVar(&value, "value", "", "the value to send in the request")
	flag.Int64Var(&expiry, "expiry", 0, "when to expire the key in unix time")
	flag.StringVar(&url, "url", ":8595", "the server's url in host:port format")
	flag.Parse()

	conn, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch {
	case operation == data.OperationGet.String() || operation == data.OperationDel.String():
		req := data.Request{
			Operation: data.Operation(operation),
			Key:       key,
		}

		if err := writeRequest(conn, req); err != nil {
			fmt.Println(err)
			return
		}

		res, err := readResponse(conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(res)
	case operation == data.OperationSet.String():
		req := data.Request{
			Operation: data.OperationSet,
			Key:       key,
			Value:     value,
		}

		if err := writeRequest(conn, req); err != nil {
			fmt.Println(err)
			return
		}

		res, err := readResponse(conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(res)
	case operation == data.OperationExp.String():
		exp := time.Unix(expiry, 0)
		req := data.Request{
			Operation: data.OperationExp,
			Key:       key,
			Expiry:    exp,
		}

		if err := writeRequest(conn, req); err != nil {
			fmt.Println(err)
			return
		}

		res, err := readResponse(conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(res)
	}
}

func writeRequest(conn net.Conn, req data.Request) error {
	reqData, err := req.Marshal()
	if err != nil {
		return err
	}

	if _, err := conn.Write(reqData); err != nil {
		return err
	}

	return nil
}

func readResponse(conn net.Conn) (*data.Response, error) {
	res := data.Response{}
	payload := make([]byte, 1024)
	if _, err := conn.Read(payload); err != nil {
		return nil, err
	}

	if err := res.Unmarshal(payload); err != nil {
		return nil, err
	}

	return &res, nil
}
