package metrics

import (
	"context"
	"strconv"
	"time"

	"go-micro/errors"
	"go-micro/metrics"
	"go-micro/midware"
	"go-micro/transport"
)

// Option is metrics option.
type Option func(*options)

// WithRequests with requests counter.
func WithRequests(c metrics.Counter) Option {
	return func(o *options) {
		o.requests = c
	}
}

// WithSeconds with seconds histogram.
func WithSeconds(c metrics.Observer) Option {
	return func(o *options) {
		o.seconds = c
	}
}

type options struct {
	// counter: <client/server>_requests_code_total{kind, operation, code, reason}
	requests metrics.Counter
	// histogram: <client/server>_requests_seconds_bucket{kind, operation}
	seconds metrics.Observer
}

// Server is midware server-side metrics.
func Server(opts ...Option) midware.Midware {
	op := options{}
	for _, o := range opts {
		o(&op)
	}
	return func(handler midware.Handler) midware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				code      int
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = int(se.Code)
				reason = se.Reason
			}
			if op.requests != nil {
				op.requests.With(kind, operation, strconv.Itoa(code), reason).Inc()
			}
			if op.seconds != nil {
				op.seconds.With(kind, operation).Observe(time.Since(startTime).Seconds())
			}
			return reply, err
		}
	}
}

// Client is midware client-side metrics.
func Client(opts ...Option) midware.Midware {
	op := options{}
	for _, o := range opts {
		o(&op)
	}
	return func(handler midware.Handler) midware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				code      int
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromClientContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = int(se.Code)
				reason = se.Reason
			}
			if op.requests != nil {
				op.requests.With(kind, operation, strconv.Itoa(code), reason).Inc()
			}
			if op.seconds != nil {
				op.seconds.With(kind, operation).Observe(time.Since(startTime).Seconds())
			}
			return reply, err
		}
	}
}
