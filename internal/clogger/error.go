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
			hub.CaptureException(err)
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
			hub.CaptureException(err)
		})
	} else {
		l.Warn("sentry not in context")
	}
	l.ErrorContext(ctx, "error", "error", err)
}
