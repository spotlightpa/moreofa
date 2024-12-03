package commentthan

import (
	"io"
	"net/http"

	"github.com/spotlightpa/moreofa/internal/clogger"
	"github.com/spotlightpa/moreofa/internal/errx"
)

func (app *appEnv) replyError(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clogger.FromContext(r.Context()).ErrorContext(r.Context(), "error", "error", err)
		w.WriteHeader(errx.StatusCode(err))
		io.WriteString(w, errx.UserMessage(err))
	})
}
