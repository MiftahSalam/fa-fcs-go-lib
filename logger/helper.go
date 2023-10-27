package logger

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

func ToField(key string, val interface{}) (field Field) {
	field = Field{
		Key: key,
		Val: val,
	}
	return
}

func ExtractCtx(ctx context.Context) Context {
	if ctx == nil {
		return Context{}
	}

	val, ok := ctx.Value(ctxKey).(Context)
	if !ok {
		return Context{}
	}

	return val
}

func ctxToLog(ctx context.Context) (logRecord []zap.Field) {
	ctxVal := ExtractCtx(ctx)

	logRecord = append(logRecord, zap.String(externalIDKey, ctxVal.ExternalID))
	if len(ctxVal.JourneyID) != 0 {
		logRecord = append(logRecord, zap.String(journeyIDKey, ctxVal.JourneyID))
	}
	if len(ctxVal.ChainID) != 0 {
		logRecord = append(logRecord, zap.String(chainIDKey, ctxVal.ChainID))
	}
	return
}

func appendToLog(fields ...Field) (logRecord []zap.Field) {
	for _, field := range fields {
		logRecord = append(logRecord, formatLog(field.Key, field.Val))
	}

	return
}

func formatLog(key string, msg interface{}) (logRecord zap.Field) {
	if msg == nil {
		logRecord = zap.Any(key, struct{}{})
		return
	}

	// handle string, string is cannot be masked, just write it
	// but try to parse as json object if possible
	if str, ok := msg.(string); ok {
		var data interface{}
		if _err := json.Unmarshal([]byte(str), &data); _err != nil {
			logRecord = zap.String(key, str)
			return
		}

		logRecord = zap.Any(key, data)
		return
	}

	// not masked since it failed to convert to reflect.Value above
	logRecord = zap.Any(key, msg)
	return
}

func createLogger(stdout bool, level int8, location string, age int) *zapLogger {
	var opt = make([]ZapOption, 0)
	if stdout {
		opt = append(opt, WithStdout())
	} else {
		opt = append(opt, WithFileOutput(location, age))
	}

	opt = append(opt, WithLevel(level))

	log, err := NewZapLogger(opt...)
	if err != nil {
		panic(fmt.Errorf("init logger with error: %w", err))
	}

	return log
}
