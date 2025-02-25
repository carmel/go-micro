package midware

import (
	"context"
)

// Handler defines the handler invoked by Midware.
type Handler func(ctx context.Context, req interface{}) (interface{}, error)

// Midware is HTTP/gRPC transport midware.
type Midware func(Handler) Handler

// Chain returns a Midware that specifies the chained handler for endpoint.
func Chain(m ...Midware) Midware {
	return func(next Handler) Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}
