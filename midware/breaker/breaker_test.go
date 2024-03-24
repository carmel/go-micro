package breaker

import (
	"context"
	"errors"
	"testing"

	mserrors "go-micro/errors"
	"go-micro/pkg/container"
	"go-micro/transport"
)

type transportMock struct {
	kind      transport.Kind
	endpoint  string
	operation string
}

type circuitBreakerMock struct {
	err error
}

func (tr *transportMock) Kind() transport.Kind {
	return tr.kind
}

func (tr *transportMock) Endpoint() string {
	return tr.endpoint
}

func (tr *transportMock) Operation() string {
	return tr.operation
}

func (tr *transportMock) RequestHeader() transport.Header {
	return nil
}

func (tr *transportMock) ReplyHeader() transport.Header {
	return nil
}

func (c *circuitBreakerMock) Allow() error { return c.err }
func (c *circuitBreakerMock) MarkSuccess() {}
func (c *circuitBreakerMock) MarkFailed()  {}

func Test_WithGroup(t *testing.T) {
	o := options{
		group: container.NewGroup(func() interface{} {
			return ""
		}),
	}

	WithGroup(nil)(&o)
	if o.group != nil {
		t.Error("The group property must be updated to nil.")
	}
}

func TestServer(_ *testing.T) {
	nextValid := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "Hello valid", nil
	}
	nextInvalid := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, mserrors.InternalServer("", "")
	}

	ctx := transport.NewClientContext(context.Background(), &transportMock{})

	_, _ = Client(func(o *options) {
		o.group = container.NewGroup(func() interface{} {
			return &circuitBreakerMock{err: errors.New("circuitbreaker error")}
		})
	})(nextValid)(ctx, nil)

	_, _ = Client(func(_ *options) {})(nextValid)(ctx, nil)

	_, _ = Client(func(_ *options) {})(nextInvalid)(ctx, nil)
}
