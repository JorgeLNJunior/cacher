package main

import (
	"bytes"
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JorgeLNJunior/cacher/pkg/data"
	levellog "github.com/JorgeLNJunior/cacher/pkg/logger"
)

const maxChunckSize = 4096

// Listen starts a tcp server.
func (app *application) Listen() error {
	listener, err := net.Listen("tcp", app.config.address)
	if err != nil {
		return err
	}
	defer listener.Close()

	app.logger.Info("tcp server is listening", levellog.Args{"addr": app.config.address})

	shutdownErr := make(chan error)
	go func() {
		exitChan := make(chan os.Signal, 1)
		signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
		<-exitChan

		app.logger.Info("started shutting down the server", nil)

		c := make(chan int)
		go func() {
			defer close(c)
			app.logger.Info("waiting for open connections before shutting down the server", nil)
			app.connectionGroup.Wait()
		}()

		select {
		case <-c:
			app.logger.Info("the server has successfully shutdown", nil)
			shutdownErr <- nil
		case <-time.After(time.Second * 5):
			app.logger.Error("the server has timed out while closing", nil)
			shutdownErr <- errors.New("timeout")
		}
	}()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					break // stop processing connections if the server is closed
				}

				app.logger.Error("error accepting a tcp connection", levellog.Args{"err": err.Error()})
				continue
			}

			go app.handleConnection(conn)
		}
	}()

	err = <-shutdownErr
	if err != nil {
		app.logger.Error("error shuting down the server", levellog.Args{"err": err.Error()})
	}

	if app.config.persist {
		app.logger.Info("started persisting the data on disk", nil)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := app.persistanceStorage.Persist(ctx, app.storage); err != nil {
			app.logger.Error("error persisting the data on disk", levellog.Args{"err": err.Error()})
			return err
		}
		app.logger.Info("the data has been successfully persisted", nil)
	}

	return nil
}

func (app *application) handleConnection(conn net.Conn) {
	defer app.connectionGroup.Done()
	defer conn.Close()
	defer func() {
		if r := recover(); r != nil {
			app.logger.Error("recovered from a panic while processing a request", nil)
		}
	}()

	app.connectionGroup.Add(1)

	if err := conn.SetWriteDeadline(time.Now().Add(time.Second * 5)); err != nil {
		app.logger.Error("error setting write timeout", levellog.Args{"err": err.Error()})
		return
	}

	buffer := bytes.NewBuffer(nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := app.readDataCtx(ctx, conn, buffer); err != nil {
		app.logger.Error("error reading data from a connection", levellog.Args{"err": err.Error()})
		app.errorResponse(conn, err)
		return
	}

	req := data.Request{}
	if err := req.Unmarshal(buffer.Bytes()); err != nil {
		app.errorResponse(conn, err)
		return
	}

	if req.Operation == data.OperationGet {
		value, ok := app.storage.Get(req.Key)
		if !ok {
			app.errorResponse(conn, errors.New("key not found"))
			return
		}
		app.okResponse(conn, value)
		return
	}
	if req.Operation == data.OperationSet {
		app.storage.Set(req.Key, req.Value)
		app.okResponse(conn, "the value has been inserted successfully")
		return
	}
	if req.Operation == data.OperationDel {
		app.storage.Delete(req.Key)
		app.okResponse(conn, "the value has been deleted successfully")
		return
	}
	if req.Operation == data.OperationExp {
		app.storage.ExpireAt(req.Key, req.Expiry)
		app.okResponse(conn, "the expiry has been set successfully")
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
	res := data.NewResponse(data.ResponseStatusError, err.Error())
	app.genericResponse(conn, res)
}

func (app *application) okResponse(conn net.Conn, message string) {
	res := data.NewResponse(data.ResponseStatusOK, message)
	app.genericResponse(conn, res)
}

func (app *application) genericResponse(conn net.Conn, res data.Response) {
	data, err := res.Marshal()
	if err != nil {
		app.logger.Error("error parsing the response", levellog.Args{"err": err.Error()})
		return
	}

	if _, err := conn.Write(data); err != nil {
		app.logger.Error("error writing data to a connection", levellog.Args{"err": err.Error()})
	}
}
