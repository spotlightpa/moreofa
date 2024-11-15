package commentthan

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strconv"

	"github.com/carlmjohnson/flagx"
	"github.com/carlmjohnson/flagx/lazyio"
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
	if err = app.Exec(); err != nil {
		clogger.Logger.Error("runtime error", "error", err)
	}
	return err
}

func (app *appEnv) ParseArgs(args []string) error {
	fl := flag.NewFlagSet(AppName, flag.ContinueOnError)
	src := lazyio.FileOrURL(lazyio.StdIO, nil)
	app.src = src
	fl.Var(src, "src", "source file or URL")
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
	if err := flagx.ParseEnv(fl, AppName); err != nil {
		return err
	}
	return nil
}

type appEnv struct {
	src io.ReadCloser
}

func (app *appEnv) Exec() (err error) {
	clogger.Logger.Info("starting")
	defer func() { clogger.Logger.Info("done") }()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	return ctx.Err()
}
