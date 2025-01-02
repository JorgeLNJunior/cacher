package main

import (
	"bytes"
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

	if err := conn.SetReadDeadline(time.Now().Add(time.Second * 15)); err != nil {
		app.logger.Printf("error setting read timeout: %s", err.Error())
		return
	}
	if err := conn.SetWriteDeadline(time.Now().Add(time.Second * 15)); err != nil {
		app.logger.Printf("error setting write timeout: %s", err.Error())
		return
	}

	buffer := bytes.NewBuffer(nil)
	if err := app.readData(conn, buffer); err != nil {
		app.logger.Printf("error reading data from a connection: %s", err.Error())

		_, err := conn.Write([]byte("error: cannot read data"))
		if err != nil {
			app.logger.Printf("error writing data to a connection: %s", err.Error())
			return
		}
	}

	req := Request{
		Operation: OperationGet,
		Key:       "foo",
		Value:     []byte("bar"),
	}

	if _, err := conn.Write([]byte(req.String())); err != nil {
		app.logger.Printf("error writing data to a connection: %s", err.Error())
		return
	}
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
