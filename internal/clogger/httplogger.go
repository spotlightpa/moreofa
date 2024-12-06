package clogger

import (
	"net/http"
	"time"

	"github.com/carlmjohnson/requests"
)

var HTTPTransport http.RoundTripper

func init() {
	HTTPTransport = requests.LogTransport(http.DefaultTransport, logReq)
	http.DefaultTransport = HTTPTransport
	http.DefaultClient.Timeout = 5 * time.Second
}

func logReq(req *http.Request, res *http.Response, err error, duration time.Duration) {
	level := SpeedThreshold(duration, 500*time.Millisecond, 1*time.Second)
	FromContext(req.Context()).
		InfoContext(req.Context(), "RoundTrip",
			"req_method", req.Method,
			"req_host", req.Host,
			"res_status", res.StatusCode,
			"res_time_class", level,
			"duration", duration,
		)
}
