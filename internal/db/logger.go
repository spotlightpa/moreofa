package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spotlightpa/moreofa/internal/clogger"
	"github.com/spotlightpa/moreofa/internal/stringx"
)

func Log(db DBTX) DBTX {
	if l, ok := db.(logger); ok {
		return l
	}
	return logger{db}
}

type logger struct {
	db DBTX
}

func (l logger) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	start := time.Now()
	s, err := l.db.PrepareContext(ctx, query)
	err = l.log(ctx, "Exec", time.Since(start), err)
	return s, err
}

func (l logger) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	t, err := l.db.ExecContext(ctx, query, args...)
	err = l.log(ctx, "Exec", time.Since(start), err)
	return t, err
}

func (l logger) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := l.db.QueryContext(ctx, query, args...)
	err = l.log(ctx, "Query", time.Since(start), err)
	return rows, err
}

func (l logger) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := l.db.QueryRowContext(ctx, query, args...)
	_ = l.log(ctx, "QueryRow", time.Since(start), nil)
	return row
}

func (l logger) log(ctx context.Context, kind string, d time.Duration, err error) error {
	pc, file, line, ok := runtime.Caller(2)
	prefix := "unknown function"
	if ok {
		f := runtime.FuncForPC(pc)
		file = filepath.Base(file)
		_, name, _ := stringx.LastCut(f.Name(), ".")
		prefix = fmt.Sprintf("%s(%s:%d)", name, file, line)
	}
	level := clogger.LevelThreshold(d, 200*time.Millisecond, 1*time.Second)
	if err != nil {
		level = slog.LevelError
		err = fmt.Errorf("%s: %w", prefix, err)
	}
	clogger.FromContext(ctx).
		Log(ctx, level, "DBTX",
			"kind", kind,
			"query", prefix,
			"duration", d,
		)
	return err
}
