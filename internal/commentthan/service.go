package commentthan

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/earthboundkid/versioninfo/v2"
	"github.com/getsentry/sentry-go"
	"github.com/spotlightpa/moreofa/internal/clogger"
	"github.com/spotlightpa/moreofa/internal/db"
)

type service struct {
	db *sql.DB
	q  *db.Queries
}

func (app *appEnv) configureService() error {
	if app.sentryDSN == "" {
		clogger.UseDevLogger()
		clogger.Logger.Warn("configureService", "Sentry-enabled", false)
	} else {
		clogger.UseProdLogger()
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:        app.sentryDSN,
			Release:    versioninfo.Revision,
			ServerName: os.Getenv("FLY_MACHINE_ID"),
		}); err != nil {
			clogger.LogErr(context.Background(), err)
		} else {
			clogger.Logger.Info("configureService", "Sentry-enabled", true)
		}
	}
	if err := db.Migrate(app.dbname); err != nil {
		return err
	}
	dbase, err := db.Open(app.dbname)
	if err != nil {
		return err
	}
	app.svc = &service{
		db: dbase,
		q:  db.New(db.Log(dbase)),
	}
	return nil
}

func (app *appEnv) closeService() {
	if err := app.svc.db.Close(); err != nil {
		clogger.Logger.Error("closeService", "error", err)
	}
	sentry.Flush(5 * time.Second)
}
