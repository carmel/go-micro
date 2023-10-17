package tracing

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	config "github.com/carmel/go-micro/gateway/api/config/v1"
	v1 "github.com/carmel/go-micro/gateway/api/midware/tracing/v1"
	"github.com/carmel/go-micro/gateway/midware"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestTracer(t *testing.T) {
	cfg, err := anypb.New(&v1.Tracing{
		HttpEndpoint: "127.0.0.1:4318",
	})
	if err != nil {
		t.Fatal(err)
	}

	next := midware.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			Body: io.NopCloser(bytes.NewBufferString("Hello Kratos")),
		}, nil
	})

	m, err := Midware(&config.Midware{
		Options: cfg,
	})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/api/v1/hello", bytes.NewBufferString("test"))
	_, err = m(next).RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
}
