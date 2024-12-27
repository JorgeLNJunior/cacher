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

	var received int
	buffer := bytes.NewBuffer(nil)

	// see: https://mostafa.dev/why-do-tcp-connections-in-go-get-stuck-reading-large-amounts-of-data-f490a26a605e
	for {
		chunck := make([]byte, maxChunckSize)

		read, err := conn.Read(chunck)
		if err != nil {
			app.logger.Printf("error reading data from a connection: %s", err.Error())
			return
		}
		received += read

		if _, err := buffer.Write(chunck[:read]); err != nil {
			app.logger.Printf("error writing data from a connection to the buffer: %s", err.Error())
			return
		}

		if read == 0 || read < maxChunckSize {
			break
		}
	}

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		app.logger.Printf("error writing to a connection: %s", err.Error())
	}
}
