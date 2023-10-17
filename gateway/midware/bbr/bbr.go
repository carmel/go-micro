package bbr

import (
	"bytes"
	"io"
	"net/http"

	config "github.com/carmel/go-micro/gateway/api/config/v1"
	"github.com/carmel/go-micro/gateway/midware"
	"github.com/carmel/go-micro/pkg/ratelimit"
	"github.com/carmel/go-micro/pkg/ratelimit/bbr"
)

var _nopBody = io.NopCloser(&bytes.Buffer{})

func init() {
	midware.Register("bbr", Midware)
}

func Midware(c *config.Midware) (midware.Midware, error) {
	limiter := bbr.NewLimiter() //use default settings
	return func(next http.RoundTripper) http.RoundTripper {
		return midware.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			done, err := limiter.Allow()
			if err != nil {
				return &http.Response{
					Status:     http.StatusText(http.StatusTooManyRequests),
					StatusCode: http.StatusTooManyRequests,
					Body:       _nopBody,
				}, nil
			}
			resp, err := next.RoundTrip(req)
			done(ratelimit.DoneInfo{Err: err})
			return resp, err
		})
	}, nil
}
