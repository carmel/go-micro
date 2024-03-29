package breaker

import (
	"context"

	"go-micro/pkg/breaker"
	"go-micro/pkg/breaker/sre"

	"go-micro/errors"
	"go-micro/midware"
	"go-micro/pkg/container"
	"go-micro/transport"
)

// ErrNotAllowed is request failed due to circuit breaker triggered.
var ErrNotAllowed = errors.New(503, "CIRCUITBREAKER", "request failed due to circuit breaker triggered")

// Option is circuit breaker option.
type Option func(*options)

// WithGroup with circuit breaker group.
// NOTE: implements generics breaker.CircuitBreaker
func WithGroup(g *container.Group) Option {
	return func(o *options) {
		o.group = g
	}
}

// WithCircuitBreaker with circuit breaker genFunc.
func WithCircuitBreaker(genBreakerFunc func() breaker.CircuitBreaker) Option {
	return func(o *options) {
		o.group = container.NewGroup(func() interface{} {
			return genBreakerFunc()
		})
	}
}

type options struct {
	group *container.Group
}

// Client breaker midware will return errBreakerTriggered when the circuit
// breaker is triggered and the request is rejected directly.
func Client(opts ...Option) midware.Midware {
	opt := &options{
		group: container.NewGroup(func() interface{} {
			return sre.NewBreaker()
		}),
	}
	for _, o := range opts {
		o(opt)
	}
	return func(handler midware.Handler) midware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			info, _ := transport.FromClientContext(ctx)
			breaker := opt.group.Get(info.Operation()).(breaker.CircuitBreaker)
			if err := breaker.Allow(); err != nil {
				// rejected
				// NOTE: when client reject requests locally,
				// continue to add counter let the drop ratio higher.
				breaker.MarkFailed()
				return nil, ErrNotAllowed
			}
			// allowed
			reply, err := handler(ctx, req)
			if err != nil && (errors.IsInternalServer(err) || errors.IsServiceUnavailable(err) || errors.IsGatewayTimeout(err)) {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
			return reply, err
		}
	}
}
