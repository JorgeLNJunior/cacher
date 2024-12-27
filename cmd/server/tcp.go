package main

import (
	"bytes"
	"net"
	"time"
)

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

		if err := conn.SetReadDeadline(time.Now().Add(time.Second * 15)); err != nil {
			app.logger.Printf("error setting read timeout: %s", err.Error())
			conn.Close()
			continue
		}
		if err := conn.SetWriteDeadline(time.Now().Add(time.Second * 15)); err != nil {
			app.logger.Printf("error setting write timeout: %s", err.Error())
			conn.Close()
			continue
		}

		buffer := bytes.NewBuffer(nil)
		chunck := make([]byte, app.config.payloadSizeLimit)

		read, err := conn.Read(chunck)
		if err != nil {
			app.logger.Printf("error reading data from a connection: %s", err.Error())
			conn.Close()
			continue
		}

		if _, err := buffer.Read(chunck[:read]); err != nil {
			app.logger.Printf("error reading data from a connection: %s", err.Error())
			continue
		}

		if _, err := conn.Write(buffer.Bytes()); err != nil {
			app.logger.Printf("error writing to a connection: %s", err.Error())
			conn.Close()
			continue
		}

		conn.Close()
	}
}
