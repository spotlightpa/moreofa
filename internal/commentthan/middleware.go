package commentthan

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/earthboundkid/mid"
	"github.com/earthboundkid/versioninfo/v2"
	"github.com/spotlightpa/moreofa/internal/errx"
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

func (svc *service) userMiddleware(h http.Handler) http.Handler {
	return svc.sessionManager.LoadAndSave(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := svc.sessionManager.Get(r.Context(), "user").(User)
			if !ok {
				user.Name = "Anonymous"
			}
			r2 := ReqWithUser(r, user)
			h.ServeHTTP(w, r2)
		}),
	)
}

func (rr *router) roleMiddleware(role string) mid.Middleware {
	return func(h http.Handler) http.Handler {
		return mid.Controller(func(w http.ResponseWriter, r *http.Request) http.Handler {
			if user := UserFromReq(r); !slices.Contains(user.Role, role) {
				return rr.replyHTMLErr(errx.E{
					S: http.StatusForbidden,
					M: "Not authorized",
					E: fmt.Errorf("need role %q had roles %q", role, user.Role),
				})
			}
			return h
		})
	}
}
