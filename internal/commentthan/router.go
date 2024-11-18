package commentthan

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/earthboundkid/mid"
	"github.com/spotlightpa/moreofa/internal/clogger"
	"github.com/spotlightpa/moreofa/static"
)

func (app *appEnv) router() http.Handler {
	srv := http.NewServeMux()
	srv.Handle("GET /", app.notFound())
	srv.Handle("GET /api/healthcheck", app.healthCheck())

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
		clogger.Middleware,
		maxBytesMiddleware(fiveMB),
		timeoutMiddleware(10 * time.Second),
		versionMiddleware,
	}
	return baseMW.Handler(srv)
}
