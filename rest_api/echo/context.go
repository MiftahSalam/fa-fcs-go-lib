package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/MiftahSalam/fa-fcs-go-lib/context"
	"github.com/MiftahSalam/fa-fcs-go-lib/errors"
	"github.com/MiftahSalam/fa-fcs-go-lib/logger"
)

const (
	CustomContext = "CustomContext"
	Session       = "Session"
)

type (
	ApplicationContext struct {
		echo.Context
		CustomContext *context.CustomContext
	}

	Success struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	Failed struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Error   string      `json:"error"`
		Data    interface{} `json:"data"`
	}
)

func (sc *ApplicationContext) Success(data interface{}) error {
	hc := http.StatusOK
	if data == nil {
		data = struct{}{}
	}

	res := Success{
		Code:    "00",
		Message: "success",
		Data:    data,
	}

	sc.CustomContext.Response = res
	sc.CustomContext.Info("Outgoing",
		logger.ToField("rt", sc.CustomContext.ResponseTime()),
		logger.ToField("response", res),
		logger.ToField("http_code", hc))
	sc.CustomContext.Summary()
	sc.CustomContext.ErrorCode = res.Code
	sc.CustomContext.ErrorMessage = res.Message

	return sc.JSON(hc, res)
}

func (sc *ApplicationContext) SuccessWithMeta(data, meta interface{}) error {
	hc := http.StatusOK
	res := Success{
		Code:    "00",
		Message: "success",
		Data:    data,
	}

	sc.CustomContext.Response = res
	sc.CustomContext.Info("Outgoing",
		logger.ToField("rt", sc.CustomContext.ResponseTime()),
		logger.ToField("response", res),
		logger.ToField("http_code", hc))
	sc.CustomContext.Summary()
	sc.CustomContext.ErrorCode = res.Code
	sc.CustomContext.ErrorMessage = res.Message

	return sc.JSON(hc, res)
}

func (sc *ApplicationContext) Fail(err error) error {
	return sc.FailWithData(err, nil)
}

func (sc *ApplicationContext) FailWithData(err error, data interface{}) error {
	var (
		ed = errors.ExtractError(err)
	)

	if data == nil {
		data = struct{}{}
	}

	res := Failed{
		Code:    ed.Code,
		Message: ed.Message,
		Error:   ed.FullMessage,
		Data:    data,
	}

	sc.CustomContext.Response = res
	sc.CustomContext.Info("Outgoing",
		logger.ToField("rt", sc.CustomContext.ResponseTime()),
		logger.ToField("response", res),
		logger.ToField("http_code", ed.HttpCode))
	sc.CustomContext.Summary()
	sc.CustomContext.ErrorCode = res.Code
	sc.CustomContext.ErrorMessage = res.Message

	return sc.JSON(ed.HttpCode, res)
}

func (sc *ApplicationContext) Raw(hc int, data interface{}) error {
	if data == nil {
		data = struct{}{}
	}

	sc.CustomContext.Response = data
	sc.CustomContext.Info("Outgoing",
		logger.ToField("rt", sc.CustomContext.ResponseTime()),
		logger.ToField("response", data),
		logger.ToField("http_code", hc))
	sc.CustomContext.Summary()

	return sc.JSON(hc, data)
}

func (sc *ApplicationContext) AddCustomContext(rc *context.CustomContext) *ApplicationContext {
	sc.Set(CustomContext, rc)
	sc.CustomContext = rc
	return sc
}

func ParseApplicationContext(c echo.Context) *ApplicationContext {
	var (
		nc  = c.Get(CustomContext)
		ctx *context.CustomContext
	)

	// request context is mandatory on application context
	// force casting
	ctx = nc.(*context.CustomContext)

	return &ApplicationContext{Context: c, CustomContext: ctx}
}

func NewEmptyApplicationContext(parent echo.Context) *ApplicationContext {
	return &ApplicationContext{parent, nil}
}

func NewApplicationContext(parent echo.Context) (*ApplicationContext, error) {
	pctx, ok := parent.(*ApplicationContext)
	if !ok {
		return nil, errors.ErrSession
	}

	return pctx, nil
}
