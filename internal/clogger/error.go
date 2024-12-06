package clogger

import (
	"context"
	"net/http"

	"github.com/getsentry/sentry-go"
)

func LogErr(ctx context.Context, err error) {
	l := FromContext(ctx)
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			e := hub.CaptureException(err)
			l.InfoContext(ctx, "sentry", "id", *e)
		})
	} else {
		l.Warn("sentry not in context")
	}
	l.ErrorContext(ctx, "error", "error", err)
}

func LogRequestErr(r *http.Request, err error) {
	ctx := r.Context()
	l := FromContext(ctx)
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetRequest(r)
			e := hub.CaptureException(err)
			l.InfoContext(ctx, "sentry", "id", *e)
		})
	} else {
		l.Warn("sentry not in context")
	}
	l.ErrorContext(ctx, "error", "error", err)
}
