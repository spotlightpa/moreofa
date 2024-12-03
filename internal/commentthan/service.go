package commentthan

import (
	"database/sql"

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
	} else {
		clogger.UseProdLogger()
		sentry.Init(sentry.ClientOptions{
			Dsn:     app.sentryDSN,
			Release: versioninfo.Revision,
		})
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
}
