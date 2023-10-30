package errors

import (
	"github.com/joomcode/errorx"
)

type (
	ErrorDescription struct {
		Code                         string
		HttpCode                     int
		Message, FullMessage, Source string
	}
)

var (
	ErrCodeProperty     = errorx.RegisterProperty("code")
	ErrHttpCodeProperty = errorx.RegisterProperty("httpcode")
	ErrSourceProperty   = errorx.RegisterProperty("source")
	ErrMessage          = errorx.RegisterProperty("message")
	ErrNamespace        = errorx.NewNamespace("fa-fcs")
	ErrBase             = errorx.NewType(ErrNamespace, "base")

	ErrSessionHeader = ErrBase.New("Authorization header is empty").WithProperty(ErrCodeProperty, "401").WithProperty(ErrHttpCodeProperty, 401)

	// - session
	ErrExpiredSession = ErrBase.New("session is already expired").WithProperty(ErrCodeProperty, "1000").WithProperty(ErrHttpCodeProperty, 401)
	ErrSession        = ErrBase.New("unauthorized").WithProperty(ErrCodeProperty, "1002").WithProperty(ErrHttpCodeProperty, 401)

	// - json
	ErrJsonMarshal   = ErrBase.New("failed marshal to json").WithProperty(ErrCodeProperty, "1003").WithProperty(ErrHttpCodeProperty, 400)
	ErrJsonUnmarshal = ErrBase.New("failed unmarshal from json").WithProperty(ErrCodeProperty, "1003").WithProperty(ErrHttpCodeProperty, 400)

	// - validation
	ErrValidation = ErrBase.New("failed to validate request body").WithProperty(ErrCodeProperty, "1004").WithProperty(ErrHttpCodeProperty, 400)

	// - connection
	ErrConnection = ErrBase.New("failed connection").WithProperty(ErrCodeProperty, "1005").WithProperty(ErrHttpCodeProperty, 500)
)

func WrapErr(err error, message string) *errorx.Error {
	return errorx.Decorate(err, message)
}

func ExtractError(err error) ErrorDescription {
	var (
		e, ok = err.(*errorx.Error)
	)

	if ok {
		if ErrNamespace.IsNamespaceOf(e.Type()) {
			code, source, httpcode := "0", "internal", 0
			c, ok := errorx.ExtractProperty(e, ErrCodeProperty)

			if ok {
				code = c.(string)
			} else {
				code = "500"
			}

			hc, ok := errorx.ExtractProperty(e, ErrHttpCodeProperty)

			if ok {
				httpcode = hc.(int)
			} else {
				httpcode = 500
			}

			s, ok := errorx.ExtractProperty(e, ErrSourceProperty)

			if ok {
				source = s.(string)
			}

			return ErrorDescription{code, httpcode, e.Message(), e.Error(), source}
		}
	}

	return ErrorDescription{
		Code:        "500",
		HttpCode:    500,
		Message:     "internal error",
		FullMessage: err.Error(),
		Source:      "internal",
	}
}

func New(hc int, code, message string) *errorx.Error {
	return ErrBase.New(message).
		WithProperty(ErrCodeProperty, code).
		WithProperty(ErrHttpCodeProperty, hc)
}

func NewWithSource(hc int, code, message, source string) *errorx.Error {
	return ErrBase.New(message).
		WithProperty(ErrCodeProperty, code).
		WithProperty(ErrHttpCodeProperty, hc).
		WithProperty(ErrSourceProperty, source)
}
