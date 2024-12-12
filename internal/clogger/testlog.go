package clogger

import (
	"bytes"
	"log/slog"
	"testing"
)

func UseTestLogger(t testing.TB) {
	opts := slog.HandlerOptions{
		Level:       Level,
		ReplaceAttr: removeTime,
	}
	logger := slog.New(slog.NewTextHandler(tWriter{t}, &opts))
	slog.SetDefault(logger)
}

type tWriter struct {
	t testing.TB
}

func (tw tWriter) Write(data []byte) (int, error) {
	tw.t.Log(string(bytes.TrimSpace(data)))
	return len(data), nil
}
