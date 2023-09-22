package binding

import (
	"net/http"
	"net/url"

	"github.com/carmel/microservices/codec"
	"github.com/carmel/microservices/codec/form"
	"github.com/carmel/microservices/errors"
)

// BindQuery bind vars parameters to target.
func BindQuery(vars url.Values, target interface{}) error {
	if err := codec.GetCodec(form.Name).Unmarshal([]byte(vars.Encode()), target); err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	return nil
}

// BindForm bind form parameters to target.
func BindForm(req *http.Request, target interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := codec.GetCodec(form.Name).Unmarshal([]byte(req.Form.Encode()), target); err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	return nil
}