package validate

import (
	"context"

	"go-micro/errors"
	"go-micro/midware"
)

type validator interface {
	Validate() error
}

// Validator is a validator midware.
func Validator() midware.Midware {
	return func(handler midware.Handler) midware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if v, ok := req.(validator); ok {
				if err := v.Validate(); err != nil {
					return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
				}
			}
			return handler(ctx, req)
		}
	}
}
