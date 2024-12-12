package commentthan

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/earthboundkid/mid"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/spotlightpa/moreofa/internal/clogger"
	"github.com/spotlightpa/moreofa/static"
)

type router struct {
	svc *service
	h   http.Handler
}

var _ http.Handler = (*router)(nil)

func (rr *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rr.h.ServeHTTP(w, r)
}

func (svc *service) router() http.Handler {
	rr := &router{svc: svc}
	srv := http.NewServeMux()
	srv.Handle("GET /", rr.notFound())
	srv.Handle("GET /api/healthcheck", rr.healthCheck())
	srv.Handle("GET /api/sentrycheck", rr.sentryCheck())
	srv.Handle("POST /comment", rr.postComment())

	fs.WalkDir(static.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		srv.HandleFunc("GET /"+path, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, static.FS, path)
		})
		return nil
	})

	const fiveMB = 5 * 1 << 20
	baseMW := mid.Stack{
		sentryhttp.New(sentryhttp.Options{}).Handle,
		clogger.Middleware,
		maxBytesMiddleware(fiveMB),
		timeoutMiddleware(10 * time.Second),
		versionMiddleware,
	}
	rr.h = baseMW.Handler(srv)
	return rr
}
