package middleware

import (
	"bytes"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	icontext "github.com/MiftahSalam/fa-fcs-go-lib/context"
	"github.com/MiftahSalam/fa-fcs-go-lib/logger"
	iecho "github.com/MiftahSalam/fa-fcs-go-lib/rest_api/echo"
)

const (
	ExternalID = "X-EXTERNAL-ID"
	JourneyID  = "X-JOURNEY-ID"
	ChainID    = "X-CHAIN-ID"
)

type (
	ContextInjectorMiddleware interface {
		Injector(next echo.HandlerFunc) echo.HandlerFunc
	}

	contextInjectorMiddleware struct {
		logger     logger.Logger
		prefixSkip []string
	}
)

func (i *contextInjectorMiddleware) Injector(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			tid = c.Request().Header.Get(ExternalID)
			jid = c.Request().Header.Get(JourneyID)
			cid = c.Request().Header.Get(ChainID)
		)
		if len(tid) == 0 {
			_uuid, _ := uuid.NewRandom()
			tid = _uuid.String()
		}

		// - Set session to context
		ctx := icontext.NewCustomContext(i.logger, tid, jid, cid)
		ctx.Context = c.Request().Context()
		ctx.Header = c.Request().Header
		ctx.URI = c.Request().URL.String()

		c.Set(iecho.CustomContext, ctx)

		// print request time
		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request().Body)
			// Restore the io.ReadCloser to its original state
			c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		ctx.Request = string(bodyBytes)
		if !i.skipper(c) {
			ctx.Info("Incoming",
				logger.ToField("url", c.Request().URL.String()),
				logger.ToField("header", ctx.Header),
				logger.ToField("request", ctx.Request))
		}

		return h(c)
	}
}

func (i *contextInjectorMiddleware) skipper(c echo.Context) (skip bool) {
	url := c.Request().URL.String()
	if url == "/" {
		skip = true
		return
	}

	for _, urlSkip := range i.prefixSkip {
		if strings.HasPrefix(url, urlSkip) {
			skip = true
			return
		}
	}

	return
}

func NewContextInjectorMiddleware(logger logger.Logger, prefixSkip ...string) (ContextInjectorMiddleware, error) {
	return &contextInjectorMiddleware{logger: logger, prefixSkip: prefixSkip}, nil
}
