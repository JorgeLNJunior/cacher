package logger

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type LogLevel int8

const (
	LevelInfo LogLevel = iota
	LevelWarn
	LevelError
	LevelFatal
)

func (l LogLevel) String() string {
	switch {
	case l == LevelInfo:
		return "INFO"
	case l == LevelWarn:
		return "WARN"
	case l == LevelError:
		return "ERROR"
	case l == LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	minLevel LogLevel
	out      io.Writer
	mu       sync.Mutex
}

func NewLogger(level LogLevel, out io.Writer) *Logger {
	return &Logger{
		minLevel: level,
		out:      out,
	}
}

func (l *Logger) Info(message string, args map[string]string) {
	l.print(LevelInfo, message, args)
}

func (l *Logger) Warn(message string, args map[string]string) {
	l.print(LevelWarn, message, args)
}

func (l *Logger) Error(message string, args map[string]string) {
	l.print(LevelError, message, args)
}

func (l *Logger) Fatal(message string, args map[string]string) {
	l.print(LevelFatal, message, args)
	os.Exit(1)
}

func (l *Logger) print(level LogLevel, message string, args map[string]string) {
	if level > l.minLevel {
		return
	}

	data := struct {
		Level   LogLevel          `json:"level"`
		Time    string            `json:"time"`
		Message string            `json:"message"`
		Args    map[string]string `json:"args,omitempty"`
		Trace   string            `json:"trace,omitempty"`
	}{
		Level:   level,
		Time:    time.Now().UTC().Format(time.RFC3339),
		Message: message,
		Args:    args,
	}

	if level >= LevelError {
		data.Trace = string(debug.Stack())
	}

	line, err := json.Marshal(data)
	if err != nil {
		line = []byte(LevelError.String() + ": error marshaling the log message: " + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	_, _ = l.out.Write(append(line, '\n'))
}
