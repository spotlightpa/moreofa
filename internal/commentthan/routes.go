package commentthan

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	"github.com/earthboundkid/mid"
	"github.com/spotlightpa/moreofa/internal/db"
	"github.com/spotlightpa/moreofa/internal/errx"
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

func (app *appEnv) postComment() mid.Controller {
	return func(w http.ResponseWriter, r *http.Request) http.Handler {
		if err := r.ParseForm(); err != nil {
			return app.replyError(errx.E{S: http.StatusBadRequest, E: err})
		}
		name := r.PostForm.Get("name")
		contact := r.PostForm.Get("contact")
		message := r.PostForm.Get("message")
		var comment db.Comment
		err := db.Tx(r.Context(), app.srv.db, &sql.TxOptions{ReadOnly: false}, func(qtx *db.Queries) error {
			var err error
			comment, err = qtx.CreateComment(r.Context(), db.CreateCommentParams{
				Name:      name,
				Contact:   contact,
				Message:   message,
				Ip:        r.RemoteAddr,
				UserAgent: r.UserAgent(),
				Referrer:  r.Referer(),
			})
			return err
		})
		if err != nil {
			return app.replyError(err)
		}
		io.WriteString(w, fmt.Sprint(comment.ID))
		return nil
	}
}
