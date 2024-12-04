package commentthan

import (
	"bytes"
	"net/http"

	"github.com/spotlightpa/moreofa/internal/clogger"
	"github.com/spotlightpa/moreofa/internal/errx"
	"github.com/spotlightpa/moreofa/layouts"
)

func (app *appEnv) replyHTMLErr(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clogger.LogErr(r.Context(), err)
		code := errx.StatusCode(err)
		var buf bytes.Buffer
		if err := layouts.Error.Execute(&buf, struct {
			Status     string
			StatusCode int
			Message    string
		}{
			Status:     http.StatusText(code),
			StatusCode: code,
			Message:    errx.UserMessage(err),
		}); err != nil {
			clogger.LogErr(r.Context(), err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(code)
		if _, err := buf.WriteTo(w); err != nil {
			clogger.LogErr(r.Context(), err)
			return
		}
	})
}
