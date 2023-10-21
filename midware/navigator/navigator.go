package navigator

import (
	"context"
	"regexp"
	"strings"

	"go-micro/midware"
	"go-micro/transport"
)

type (
	transporter func(ctx context.Context) (transport.Transporter, bool)
	MatchFunc   func(ctx context.Context, operation string) bool
)

var (
	// serverTransporter is get server transport.Transporter from ctx
	serverTransporter transporter = func(ctx context.Context) (transport.Transporter, bool) {
		return transport.FromServerContext(ctx)
	}
	// clientTransporter is get client transport.Transporter from ctx
	clientTransporter transporter = func(ctx context.Context) (transport.Transporter, bool) {
		return transport.FromClientContext(ctx)
	}
)

// Builder is a navigator builder
type Builder struct {
	match  MatchFunc
	prefix []string
	regex  []string
	path   []string
	ms     []midware.Midware
	client bool
}

// Server navigator midware
func Server(ms ...midware.Midware) *Builder {
	return &Builder{ms: ms}
}

// Client navigator midware
func Client(ms ...midware.Midware) *Builder {
	return &Builder{client: true, ms: ms}
}

// Prefix is with Builder's prefix
func (b *Builder) Prefix(prefix ...string) *Builder {
	b.prefix = prefix
	return b
}

// Regex is with Builder's regex
func (b *Builder) Regex(regex ...string) *Builder {
	b.regex = regex
	return b
}

// Path is with Builder's path
func (b *Builder) Path(path ...string) *Builder {
	b.path = path
	return b
}

// Match is with Builder's match
func (b *Builder) Match(fn MatchFunc) *Builder {
	b.match = fn
	return b
}

// Build is Builder's Build, for example: Server().Path(m1,m2).Build()
func (b *Builder) Build() midware.Midware {
	var transporter func(ctx context.Context) (transport.Transporter, bool)
	if b.client {
		transporter = clientTransporter
	} else {
		transporter = serverTransporter
	}
	return navigator(transporter, b.matches, b.ms...)
}

// matches is match operation compliance Builder
func (b *Builder) matches(ctx context.Context, transporter transporter) bool {
	info, ok := transporter(ctx)
	if !ok {
		return false
	}

	operation := info.Operation()
	for _, prefix := range b.prefix {
		if prefixMatch(prefix, operation) {
			return true
		}
	}
	for _, regex := range b.regex {
		if regexMatch(regex, operation) {
			return true
		}
	}
	for _, path := range b.path {
		if pathMatch(path, operation) {
			return true
		}
	}

	if b.match != nil {
		if b.match(ctx, operation) {
			return true
		}
	}

	return false
}

// navigator midware
func navigator(transporter transporter, match func(context.Context, transporter) bool, ms ...midware.Midware) midware.Midware {
	return func(handler midware.Handler) midware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if !match(ctx, transporter) {
				return handler(ctx, req)
			}
			return midware.Chain(ms...)(handler)(ctx, req)
		}
	}
}

func pathMatch(path string, operation string) bool {
	return path == operation
}

func prefixMatch(prefix string, operation string) bool {
	return strings.HasPrefix(operation, prefix)
}

func regexMatch(regex string, operation string) bool {
	r, err := regexp.Compile(regex)
	if err != nil {
		return false
	}
	return r.FindString(operation) == operation
}
