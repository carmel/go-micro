package logging

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	config "github.com/carmel/microservices/gateway/api/config/v1"
	"github.com/carmel/microservices/gateway/midware"
	"github.com/carmel/microservices/logger"
)

func init() {
	midware.Register("logging", Midware)
}

// Midware is a logging midware.
func Midware(c *config.Midware) (midware.Midware, error) {
	return func(next http.RoundTripper) http.RoundTripper {
		return midware.RoundTripperFunc(func(req *http.Request) (reply *http.Response, err error) {
			startTime := time.Now()
			reply, err = next.RoundTrip(req)
			level := logger.INFO
			code := http.StatusBadGateway
			errMsg := ""
			if err != nil {
				level = logger.ERROR
				errMsg = err.Error()
			} else {
				code = reply.StatusCode
			}
			ctx := req.Context()
			// nodes, _ := midware.RequestBackendsFromContext(ctx)
			reqOpt, _ := midware.FromRequestContext(ctx)
			logger.With(
				"source", "accesslog",
				"host", req.Host,
				"method", req.Method,
				"scheme", req.URL.Scheme,
				"path", req.URL.Path,
				"query", req.URL.RawQuery,
				"code", code,
				"latency", time.Since(startTime).Seconds(),
				"backend", strings.Join(reqOpt.Backends, ","),
				"backend_code", reqOpt.UpstreamStatusCode,
				"backend_latency", reqOpt.UpstreamResponseTime,
				"last_attempt", reqOpt.LastAttempt,
			).Log(ctx, slog.Level(level), errMsg)

			return reply, err
		})
	}, nil
}
