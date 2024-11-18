package commentthan

import (
	"io"
	"net/http"

	"github.com/earthboundkid/mid"
	"github.com/spotlightpa/moreofa/static"
)

func (app *appEnv) notFound() mid.Controller {
	return func(w http.ResponseWriter, r *http.Request) http.Handler {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFileFS(w, r, static.FS, "404.html")
		return nil
	}
}

func (app *appEnv) healthCheck() mid.Controller {
	return func(w http.ResponseWriter, r *http.Request) http.Handler {
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "OK")
		return nil
	}
}
