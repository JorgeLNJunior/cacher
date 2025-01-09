package main

import (
	"bytes"
	"context"
	"errors"
	"net"
	"time"
)

const maxChunckSize = 4096

func (app *application) Listen() error {
	listener, err := net.Listen("tcp", app.config.address)
	if err != nil {
		return err
	}
	defer listener.Close()

	app.logger.Printf("tcp server listening at %s\n", app.config.address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			app.logger.Printf("error accepting a tcp connection: %s", err.Error())
			continue
		}

		go app.handleConnection(conn)
	}
}

func (app *application) handleConnection(conn net.Conn) {
	defer conn.Close()

	if err := conn.SetWriteDeadline(time.Now().Add(time.Second * 5)); err != nil {
		app.logger.Printf("error setting write timeout: %s", err.Error())
		return
	}

	buffer := bytes.NewBuffer(nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := app.readDataCtx(ctx, conn, buffer); err != nil {
		app.logger.Printf("error reading data from a connection: %s", err.Error())

		app.errorResponse(conn, err)
		return
	}

	req := Request{}
	if err := req.Unmarshal(buffer.Bytes()); err != nil {
		app.errorResponse(conn, err)
		return
	}

	if req.Operation == OperationGet {
		value, ok := app.store.Get(req.Key)
		if !ok {
			app.errorResponse(conn, errors.New("key not found"))
			return
		}
		app.okResponse(conn, value)
		return
	}
	if req.Operation == OperationSet {
		app.store.Set(req.Key, req.Value)
		app.okResponse(conn, "the value has been inserted successfully")
		return
	}

	app.errorResponse(conn, errors.New("unknown error"))
}

func (app *application) readData(conn net.Conn, to *bytes.Buffer) error {
	var received int

	// see: https://mostafa.dev/why-do-tcp-connections-in-go-get-stuck-reading-large-amounts-of-data-f490a26a605e
	for {
		chunck := make([]byte, maxChunckSize)

		read, err := conn.Read(chunck)
		if err != nil {
			return err
		}
		received += read

		if _, err := to.Write(chunck[:read]); err != nil {
			return err
		}

		if read == 0 || read < maxChunckSize {
			break
		}
	}

	return nil
}

func (app *application) readDataCtx(ctx context.Context, conn net.Conn, to *bytes.Buffer) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := app.readData(conn, to); err != nil {
			return err
		}
	}
	return nil
}

func (app *application) errorResponse(conn net.Conn, err error) {
	res := NewResponse(ResponseStatusError, err.Error())
	app.genericResponse(conn, res)
}

func (app *application) okResponse(conn net.Conn, message string) {
	res := NewResponse(ResponseStatusOK, message)
	app.genericResponse(conn, res)
}

func (app *application) genericResponse(conn net.Conn, res Response) {
	data, err := res.Marshal()
	if err != nil {
		app.logger.Printf("error parsing the response: %s", err.Error())
		return
	}

	if _, err := conn.Write(data); err != nil {
		app.logger.Printf("error writing data to a connection: %s", err.Error())
	}
}
