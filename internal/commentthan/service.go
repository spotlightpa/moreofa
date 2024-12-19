package commentthan

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/earthboundkid/versioninfo/v2"
	"github.com/getsentry/sentry-go"
	"github.com/spotlightpa/moreofa/internal/clogger"
	"github.com/spotlightpa/moreofa/internal/db"
)

type service struct {
	db              *sql.DB
	q               *db.Queries
	cl              *http.Client
	redirectSuccess string
	sessionManager  *scs.SessionManager
}

func (app *appEnv) configureService() (*service, error) {
	if app.isLocalhost {
		clogger.UseDevLogger()
		slog.Warn("configureService", "is-localhost", true)
	} else {
		clogger.UseProdLogger()
		slog.Info("configureService", "is-localhost", false)
	}
	if app.sentryDSN == "" {
		slog.Warn("configureService", "Sentry-enabled", false)
	} else {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:        app.sentryDSN,
			Release:    versioninfo.Revision,
			ServerName: os.Getenv("FLY_MACHINE_ID"),
		}); err != nil {
			clogger.LogErr(context.Background(), err)
		} else {
			slog.Info("configureService", "Sentry-enabled", true)
		}
	}
	if err := db.Migrate(app.dbname); err != nil {
		return nil, err
	}
	dbase, err := db.Open(app.dbname)
	if err != nil {
		return nil, err
	}

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(dbase)
	sessionManager.Lifetime = 28 * 24 * time.Hour
	sessionManager.Cookie.Secure = !app.isLocalhost

	return &service{
		db: dbase,
		q:  db.New(db.Log(dbase)),
		cl: &http.Client{
			Transport: clogger.HTTPTransport,
			Timeout:   5 * time.Second,
		},
		redirectSuccess: app.redirectSuccess,
		sessionManager:  sessionManager,
	}, nil
}

func (svc *service) closeService() {
	svc.sessionManager.Store.(*sqlite3store.SQLite3Store).StopCleanup()
	if err := svc.db.Close(); err != nil {
		slog.Error("closeService", "error", err)
	}
	sentry.Flush(5 * time.Second)
}
