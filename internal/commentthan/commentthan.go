package commentthan

import (
	"cmp"
	"context"
	"flag"
	"fmt"
	"io"
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
		clogger.Logger.Error("runtime error", "error", err)
	}
	return err
}

func (app *appEnv) ParseArgs(args []string) error {
	fl := flag.NewFlagSet(AppName, flag.ContinueOnError)
	fl.StringVar(&app.addr, "port", cmp.Or(os.Getenv("PORT"), ":58448"), "")
	clogger.UseDevLogger()
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
	addr string
}

func (app *appEnv) Exec(ctx context.Context) (err error) {
	defer func() { clogger.Logger.Info("done") }()

	handler := app.router()
	srv := &http.Server{
		Addr:              app.addr,
		Handler:           handler,
		BaseContext:       func(net.Listener) context.Context { return ctx },
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       1 * time.Minute,
	}
	ch := make(chan error, 1)
	go func() {
		<-ctx.Done()
		clogger.Logger.Info("shutting down")

		shutdownCtx, stop := context.WithTimeout(context.Background(), 10*time.Second)
		defer stop()
		ch <- srv.Shutdown(shutdownCtx)
	}()
	clogger.Logger.Info("starting", "port", app.addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return <-ch
}

func (app *appEnv) router() http.Handler {
	routes := clogger.Middleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			io.WriteString(w, "Hello, World!")
		}),
	)
	return routes
}
