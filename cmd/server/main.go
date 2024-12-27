package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	const address = ":8595"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Printf("error listening at %s: %s", address, err.Error())
	}
	defer listener.Close()

	logger.Printf("tcp server listening at %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Printf("error accepting a tcp connection: %s", err.Error())
			continue
		}

		if err := conn.SetReadDeadline(time.Now().Add(time.Second * 15)); err != nil {
			logger.Printf("error setting read timeout: %s", err.Error())
			conn.Close()
			continue
		}
		if err := conn.SetWriteDeadline(time.Now().Add(time.Second * 15)); err != nil {
			logger.Printf("error setting write timeout: %s", err.Error())
			conn.Close()
			continue
		}

		buffer := bytes.NewBuffer(nil)
		chunck := make([]byte, 4096)

		read, err := conn.Read(chunck)
		if err != nil {
			logger.Printf("error reading data from a connection: %s", err.Error())
			conn.Close()
			continue
		}

		if _, err := buffer.Read(chunck[:read]); err != nil {
			switch {
			case errors.Is(err, io.EOF):
				_, _ = conn.Write([]byte("entity too large"))
			default:
				logger.Printf("error reading data from a connection: %s", err.Error())
			}
			conn.Close()
			continue
		}

		if _, err := conn.Write(buffer.Bytes()); err != nil {
			logger.Printf("error writing to a connection: %s", err.Error())
			conn.Close()
			continue
		}

		conn.Close()
	}
}
