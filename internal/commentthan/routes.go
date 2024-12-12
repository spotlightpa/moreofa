package commentthan

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/earthboundkid/mid"
	"github.com/gorilla/schema"
	"github.com/spotlightpa/moreofa/internal/clogger"
	"github.com/spotlightpa/moreofa/internal/db"
	"github.com/spotlightpa/moreofa/internal/errx"
	"github.com/spotlightpa/moreofa/static"
)

func (rr *router) notFound() mid.Controller {
	return func(w http.ResponseWriter, r *http.Request) http.Handler {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFileFS(w, r, static.FS, "404.html")
		return nil
	}
}

func (rr *router) healthCheck() mid.Controller {
	return func(w http.ResponseWriter, r *http.Request) http.Handler {
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "OK")
		return nil
	}
}

func (rr *router) sentryCheck() mid.Controller {
	return func(w http.ResponseWriter, r *http.Request) http.Handler {
		clogger.LogRequestErr(r, errors.New("sentry check"))
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "OK")
		return nil
	}
}

func (rr *router) postComment() mid.Controller {
	return func(w http.ResponseWriter, r *http.Request) http.Handler {
		if err := r.ParseForm(); err != nil {
			return rr.replyHTMLErr(errx.E{S: http.StatusBadRequest, E: err})
		}
		decoder := schema.NewDecoder()
		decoder.IgnoreUnknownKeys(true)
		var req struct {
			Name      string `schema:"name"`
			Contact   string `schema:"email"`
			Subject   string `schema:"subject"`
			CC        string `schema:"CC"`
			Message   string `schema:"comment"`
			HostPage  string `schema:"host_page"`
			Anonymous bool   `schema:"anonymous"`
			BotField  string `schema:"bot-field"`
		}
		if err := decoder.Decode(&req, r.PostForm); err != nil {
			return rr.replyHTMLErr(errx.E{S: http.StatusBadRequest, E: err})
		}
		if req.Anonymous {
			req.Message = "I wish to remain anonymous.\n\n" + req.Message
		}
		l := clogger.FromContext(r.Context())
		if err := db.Tx(r.Context(), rr.svc.db, &sql.TxOptions{ReadOnly: false}, func(qtx *db.Queries) error {
			l.InfoContext(r.Context(), "postComment", "contact", req.Contact, "req_ip", r.RemoteAddr, "req_agent", r.UserAgent())
			ip, _, _ := strings.Cut(r.RemoteAddr, ":")
			_, err := qtx.CreateComment(r.Context(), db.CreateCommentParams{
				Name:      req.Name,
				Contact:   req.Contact,
				Subject:   req.Subject,
				Cc:        req.CC,
				Message:   req.Message,
				Ip:        ip,
				UserAgent: r.UserAgent(),
				Referrer:  r.Referer(),
				HostPage:  req.HostPage,
			})
			return err
		}); err != nil {
			return rr.replyHTMLErr(err)
		}
		http.Redirect(w, r, rr.svc.redirectSuccess, http.StatusSeeOther)
		return nil
	}
}
