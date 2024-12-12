// Package clogger has the common logger
package clogger

import (
	"cmp"
	"log/slog"
	"os"
	"time"
)

func init() {
	// Prevent failure to init default logger
	l := slog.New(slog.NewTextHandler(initMe{}, nil))
	slog.SetDefault(l)
}

type initMe struct{}

func (initMe) Write([]byte) (int, error) { panic("wrote to uninitialized almlog.Logger") }

var Level = &slog.LevelVar{}

func init() {
	Level.Set(slog.LevelDebug)
}

func removeTime(groups []string, a slog.Attr) slog.Attr {
	// Fly.io already logs time
	if a.Key == slog.TimeKey && len(groups) == 0 {
		a.Key = ""
		a.Value = slog.Value{}
	}
	return a
}

func UseProdLogger() {
	opts := slog.HandlerOptions{
		Level:       Level,
		ReplaceAttr: removeTime,
	}
	logger := slog.New(slog.NewTextHandler(colorize{os.Stderr}, &opts))
	slog.SetDefault(logger)
}

func shortenTime(groups []string, a slog.Attr) slog.Attr {
	// Omit date from dev
	if a.Key == slog.TimeKey && len(groups) == 0 {
		a.Value = slog.StringValue(a.Value.Time().Format("03:04:05"))
	}
	return a
}

func UseDevLogger() {
	opts := slog.HandlerOptions{
		Level:       Level,
		ReplaceAttr: shortenTime,
	}

	logger := slog.New(slog.NewTextHandler(colorize{os.Stderr}, &opts))
	slog.SetDefault(logger)
}

func SpeedThreshold(val, warn, err time.Duration) string {
	if val >= err {
		return "error"
	}
	if val >= warn {
		return "warning"
	}
	return "ok"
}

func LevelThreshold[T cmp.Ordered](val, warn, err T) slog.Level {
	if val >= err {
		return slog.LevelError
	}
	if val >= warn {
		return slog.LevelWarn
	}
	return slog.LevelInfo
}
