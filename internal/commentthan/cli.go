package commentthan

import (
	"cmp"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/carlmjohnson/flagx"
	"github.com/earthboundkid/versioninfo/v2"
	"github.com/spotlightpa/moreofa/internal/clogger"
)

const AppName = "More of a"

func CLI(args []string) error {
	var app appEnv
	err := app.ParseArgs(args)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err = app.Exec(ctx); err != nil {
		slog.Error("runtime error", "error", err)
	}
	return err
}

func (app *appEnv) ParseArgs(args []string) error {
	fl := flag.NewFlagSet(AppName, flag.ContinueOnError)
	fl.StringVar(&app.port, "port", cmp.Or(os.Getenv("PORT"), ":58448"), "")
	fl.StringVar(&app.dbname, "db", "comments.db", "")
	fl.StringVar(&app.redirectSuccess, "redirect-success", "https://www.spotlightpa.org/contact/thanks/", "")
	fl.StringVar(&app.sentryDSN, "sentry-dsn", "", "DSN `pseudo-URL` for Sentry")
	fl.Func("level", "log level", func(s string) error {
		l, _ := strconv.Atoi(s)
		clogger.Level.Set(slog.Level(l))
		return nil
	})
	fl.Usage = func() {
		fmt.Fprintf(fl.Output(), `moreofa - %s

More of a comment server than a question

Usage:

	moreofa [options]

Options:
`, versioninfo.Version)
		fl.PrintDefaults()
	}
	if err := fl.Parse(args); err != nil {
		return err
	}
	if err := flagx.ParseEnv(fl, "MOREOFA"); err != nil {
		return err
	}
	return nil
}

type appEnv struct {
	port            string
	dbname          string
	sentryDSN       string
	svc             *service
	redirectSuccess string
}

func (app *appEnv) Exec(ctx context.Context) (err error) {
	defer func() { slog.Info("done") }()

	if err := app.configureService(); err != nil {
		return err
	}
	defer app.closeService()

	handler := app.router()
	srv := &http.Server{
		Addr:              app.port,
		Handler:           handler,
		BaseContext:       func(net.Listener) context.Context { return ctx },
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       1 * time.Minute,
	}
	ch := make(chan error, 1)
	go func() {
		<-ctx.Done()
		slog.Info("shutting down")

		shutdownCtx, stop := context.WithTimeout(context.Background(), 10*time.Second)
		defer stop()
		ch <- srv.Shutdown(shutdownCtx)
	}()
	slog.Info("starting", "port", app.port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return <-ch
}
