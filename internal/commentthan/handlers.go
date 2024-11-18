package commentthan

import (
	"net/http"

	"github.com/spotlightpa/moreofa/internal/clogger"
)

func (app *appEnv) replyError(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clogger.Logger.ErrorContext(r.Context(), "error", "error", err)
		w.WriteHeader(500)
	})
}
