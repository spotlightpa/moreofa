package db

import (
	"context"
	"database/sql"
	"log"
	"net/url"

	"github.com/spotlightpa/moreofa/internal/clogger"
	"github.com/spotlightpa/moreofa/sql/migrations"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func Open(name string) (*sql.DB, error) {
	const pragmas = "?_txlock=immediate&" +
		"_pragma=busy_timeout(5000)&" +
		"_pragma=journal_mode(WAL)&" +
		"_pragma=journal_size_limit(200000000)&" +
		"_pragma=synchronous(NORMAL)&" +
		"_pragma=foreign_keys(ON)&" +
		"_pragma=temp_store(MEMORY)&" +
		"_pragma=cache_size(-16000)"

	return sql.Open("sqlite3", name+pragmas)
}

func Tx(ctx context.Context, db *sql.DB, opts *sql.TxOptions, cb func(*Queries) error) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	defer func() {
		panicking := recover()
		if panicking == nil {
			return
		}
		tx.Rollback()
		panic(panicking)
	}()
	q := New(Log(tx))
	if err = cb(q); err != nil {
		return tx.Rollback()
	}
	return tx.Commit()
}

func Migrate(name string) error {
	db := dbmate.New(&url.URL{
		Scheme: "sqlite3",
		Opaque: name,
	})
	db.FS = migrations.FS
	db.AutoDumpSchema = false
	db.MigrationsDir = []string{"."}
	db.Log = log.Writer()
	migrations, err := db.FindMigrations()
	if err != nil {
		return err
	}
	for _, m := range migrations {
		clogger.Logger.Info("db.Migrate: found", "version", m.Version, "path", m.FilePath)
	}

	clogger.Logger.Info("db.Migrate: migrating")
	return db.CreateAndMigrate()
}
