//go:build !(linux && amd64)
// +build !linux !amd64

package server

import "go-micro/transport"

func attemptSwitchingTransport(o *Options) transport.ServerTransport {
	if o.Transport == nil {
		return transport.DefaultServerStreamTransport
	}
	return o.Transport
}
