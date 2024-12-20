package commentthan

import (
	"context"
	"net/http"
	"time"

	"github.com/earthboundkid/mid"
	"github.com/earthboundkid/versioninfo/v2"
	"golang.org/x/oauth2"
)

func maxBytesMiddleware(n int64) mid.Middleware {
	return func(h http.Handler) http.Handler {
		return http.MaxBytesHandler(h, n)
	}
}

func timeoutMiddleware(timeout time.Duration) mid.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, stop := context.WithTimeout(r.Context(), timeout)
			defer stop()
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func versionMiddleware(next http.Handler) http.Handler {
	version := versioninfo.Short()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Moreofa-Version", version)
		next.ServeHTTP(w, r)
	})
}

func (svc *service) oauthClientMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, oauth2.HTTPClient, svc.cl)
		r2 := r.WithContext(ctx)
		h.ServeHTTP(w, r2)
	})
}
