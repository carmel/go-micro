package client

import (
	"io"
	"net/http"
	"time"

	"go-micro/gateway/midware"
	"go-micro/selector"
)

type client struct {
	applier  *nodeApplier
	selector selector.Selector
}

type Client interface {
	http.RoundTripper
	io.Closer
}

func newClient(applier *nodeApplier, selector selector.Selector) *client {
	return &client{
		applier:  applier,
		selector: selector,
	}
}

func (c *client) Close() error {
	c.applier.Cancel()
	return nil
}

func (c *client) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	ctx := req.Context()
	reqOpt, _ := midware.FromRequestContext(ctx)
	filter, _ := midware.SelectorFiltersFromContext(ctx)
	n, done, err := c.selector.Select(ctx, selector.WithNodeFilter(filter...))
	if err != nil {
		return nil, err
	}
	reqOpt.CurrentNode = n

	addr := n.Address()
	reqOpt.Backends = append(reqOpt.Backends, addr)
	req.URL.Host = addr
	req.URL.Scheme = "http"
	req.RequestURI = ""
	startAt := time.Now()
	resp, err = n.(*node).client.Do(req)
	reqOpt.UpstreamResponseTime = append(reqOpt.UpstreamResponseTime, time.Since(startAt).Seconds())
	if err != nil {
		done(ctx, selector.DoneInfo{Err: err})
		reqOpt.UpstreamStatusCode = append(reqOpt.UpstreamStatusCode, 0)
		return nil, err
	}
	reqOpt.UpstreamStatusCode = append(reqOpt.UpstreamStatusCode, resp.StatusCode)
	reqOpt.DoneFunc = done
	return resp, nil
}
