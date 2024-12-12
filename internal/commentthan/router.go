package commentthan

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/dghubble/gologin/v2"
	gologingoogle "github.com/dghubble/gologin/v2/google"
	"github.com/earthboundkid/mid"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/spotlightpa/moreofa/internal/clogger"
	"github.com/spotlightpa/moreofa/static"
	"golang.org/x/oauth2"
	oauth2google "golang.org/x/oauth2/google"
)

type router struct {
	svc *service
	h   http.Handler
}

var _ http.Handler = (*router)(nil)

func (rr *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rr.h.ServeHTTP(w, r)
}

func (app *appEnv) router(svc *service) http.Handler {
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

	// Add login handlers
	redirectURL := "https://moreofa.spotlightpa.org/google-callback"
	if app.isLocalhost {
		redirectURL = fmt.Sprintf("http://localhost%s/google-callback", app.port)
	}
	oauthConf := &oauth2.Config{
		ClientID:     app.googleClientID,
		ClientSecret: app.googleClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"profile", "email"},
		Endpoint:     oauth2google.Endpoint,
	}
	cookieConf := gologin.DefaultCookieConfig
	if app.isLocalhost {
		cookieConf = gologin.DebugOnlyCookieConfig
	}

	loginRedirect := gologingoogle.LoginHandler(oauthConf, rr.googleCallbackError())
	loginRedirect = gologingoogle.StateHandler(cookieConf, loginRedirect)
	srv.Handle("GET /login/{$}", loginRedirect)

	googleCallback := gologingoogle.CallbackHandler(oauthConf, rr.googleCallback(), rr.googleCallbackError())
	googleCallback = gologingoogle.StateHandler(cookieConf, googleCallback)
	googleCallback = svc.oauthClientMiddleware(googleCallback)
	srv.Handle("GET /google-callback", googleCallback)

	rr.h = baseMW.Handler(srv)
	return rr
}
