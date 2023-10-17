package midware

import (
	"io"
	"net/http"

	configv1 "github.com/carmel/go-micro/gateway/api/config/v1"
)

// Factory is a midware factory.
type Factory func(*configv1.Midware) (Midware, error)

// Midware is handler midware.
type Midware func(http.RoundTripper) http.RoundTripper

// RoundTripperFunc is an adapter to allow the use of
// ordinary functions as HTTP RoundTripper.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip calls f(w, r).
func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type FactoryV2 func(*configv1.Midware) (MidwareV2, error)
type MidwareV2 interface {
	Process(http.RoundTripper) http.RoundTripper
	io.Closer
}

func wrapFactory(in Factory) FactoryV2 {
	return func(m *configv1.Midware) (MidwareV2, error) {
		v, err := in(m)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func (f Midware) Process(in http.RoundTripper) http.RoundTripper { return f(in) }
func (f Midware) Close() error                                   { return nil }

type withCloser struct {
	process Midware
	closer  io.Closer
}

func (w *withCloser) Process(in http.RoundTripper) http.RoundTripper { return w.process(in) }
func (w *withCloser) Close() error                                   { return w.closer.Close() }
func NewWithCloser(process Midware, closer io.Closer) MidwareV2 {
	return &withCloser{
		process: process,
		closer:  closer,
	}
}

var EmptyMiddleware = emptyMidware{}

type emptyMidware struct{}

func (emptyMidware) Process(next http.RoundTripper) http.RoundTripper { return next }
func (emptyMidware) Close() error                                     { return nil }
