//go:build !(linux && amd64)
// +build !linux !amd64

package client

import "go-micro/transport"

func attemptSwitchingTransport(o *Options) transport.ClientTransport {
	if o.Transport == nil {
		return transport.DefaultClientTransport
	}
	return o.Transport
}
