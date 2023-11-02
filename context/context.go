package context

import (
	"context"
	"time"

	"github.com/MiftahSalam/fa-fcs-go-lib/logger"
)

type CustomContext struct {
	Context        context.Context
	additionalData map[string]interface{}
	Logger         logger.Logger
	RequestTime    time.Time
	ExternalID     string
	JourneyID      string
	ChainID        string
	URI            string
	Header         interface{}
	Request        interface{}
	Response       interface{}
	ErrorCode      string
	ErrorMessage   string
}

func NewCustomContext(logger logger.Logger, xid, jid, cid string) *CustomContext {
	return &CustomContext{
		Context:        context.Background(),
		additionalData: map[string]interface{}{},
		Logger:         logger,
		RequestTime:    time.Now(),
		ExternalID:     xid,
		JourneyID:      jid,
		ChainID:        cid,
		Header:         map[string]interface{}{},
		Request:        struct{}{},
		Response:       struct{}{},
	}
}

func (c *CustomContext) Get(key string) (data interface{}, ok bool) {
	data, ok = c.additionalData[key]
	return
}

func (c *CustomContext) Put(key string, data interface{}) {
	c.additionalData[key] = data
}

func (c *CustomContext) ToContextLogger() (ctx context.Context) {
	ctxVal := logger.Context{
		ExternalID:     c.ExternalID,
		JourneyID:      c.JourneyID,
		ChainID:        c.ChainID,
		AdditionalData: c.additionalData,
	}

	ctx = logger.InjectCtx(context.Background(), ctxVal)
	return
}

func (c *CustomContext) Debug(message string, field ...logger.Field) {
	c.Logger.DebugWithCtx(c.ToContextLogger(), message, field...)
}

func (c *CustomContext) Debugf(format string, arg ...interface{}) {
	c.Logger.DebugfWithCtx(c.ToContextLogger(), format, arg...)
}

func (c *CustomContext) Info(message string, field ...logger.Field) {
	c.Logger.InfoWithCtx(c.ToContextLogger(), message, field...)
}

func (c *CustomContext) Infof(format string, arg ...interface{}) {
	c.Logger.InfofWithCtx(c.ToContextLogger(), format, arg...)
}

func (c *CustomContext) Warn(message string, field ...logger.Field) {
	c.Logger.WarnWithCtx(c.ToContextLogger(), message, field...)
}

func (c *CustomContext) Warnf(format string, arg ...interface{}) {
	c.Logger.WarnfWithCtx(c.ToContextLogger(), format, arg...)
}

func (c *CustomContext) Error(message string, field ...logger.Field) {
	c.Logger.ErrorWithCtx(c.ToContextLogger(), message, field...)
}

func (c *CustomContext) Errorf(format string, arg ...interface{}) {
	c.Logger.ErrorfWithCtx(c.ToContextLogger(), format, arg...)
}

func (c *CustomContext) Fatal(message string, field ...logger.Field) {
	c.Logger.FatalWithCtx(c.ToContextLogger(), message, field...)
}

func (c *CustomContext) Fatalf(format string, arg ...interface{}) {
	c.Logger.FatalfWithCtx(c.ToContextLogger(), format, arg...)
}

func (c *CustomContext) Summary() {
	model := logger.LogSummary{
		ExternalID:     c.ExternalID,
		JourneyID:      c.JourneyID,
		ChainID:        c.ChainID,
		RespTime:       c.ResponseTime(),
		Error:          c.ErrorMessage,
		URI:            c.URI,
		Header:         c.Header,
		Request:        c.Request,
		Response:       c.Response,
		AdditionalData: c.additionalData,
	}

	c.Logger.Summary(model)
}

func (c *CustomContext) ResponseTime() int64 {
	return time.Since(c.RequestTime).Milliseconds()
}

func (c *CustomContext) GetAdditionalData() map[string]interface{} {
	return c.additionalData
}
