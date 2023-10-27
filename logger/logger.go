package logger

import (
	"context"
	"fmt"
)

type (
	Logger interface {
		Info(...interface{})
		Infof(string, ...interface{})
		Debug(...interface{})
		Debugf(string, ...interface{})
		Error(...interface{})
		Errorf(string, ...interface{})
		Warning(...interface{})
		Warningf(string, ...interface{})
		Fatal(...interface{})
		Fatalf(string, ...interface{})
		Print(...interface{})
		Printf(string, ...interface{})
		Println(...interface{})
		Instance() interface{}

		DebugWithCtx(context.Context, string, ...Field)
		DebugfWithCtx(context.Context, string, ...interface{})
		InfoWithCtx(context.Context, string, ...Field)
		InfofWithCtx(context.Context, string, ...interface{})
		WarnWithCtx(context.Context, string, ...Field)
		WarnfWithCtx(context.Context, string, ...interface{})
		ErrorWithCtx(context.Context, string, ...Field)
		ErrorfWithCtx(context.Context, string, ...interface{})
		FatalWithCtx(context.Context, string, ...Field)
		FatalfWithCtx(context.Context, string, ...interface{})
		Summary(tdr LogSummary)
	}

	LogSummary struct {
		ExternalID     string      `json:"external_id"`
		JourneyID      string      `json:"journey_id"`
		ChainID        string      `json:"chain_id"`
		RespTime       int64       `json:"rt"`
		Error          string      `json:"error"`
		URI            string      `json:"uri"`
		Header         interface{} `json:"header"`
		Request        interface{} `json:"req"`
		Response       interface{} `json:"resp"`
		AdditionalData interface{} `json:"additional_data"`
	}

	Field struct {
		Key string
		Val interface{}
	}

	Level     string
	Formatter string

	Option struct {
		Level                       Level
		LogFilePath                 string
		Formatter                   Formatter
		MaxSize, MaxBackups, MaxAge int
		Compress                    bool
	}
)

type (
	LogOption struct {
		Stdout       bool      `json:"stdout" mapstructure:"stdout"`
		FileLocation string    `json:"file_location" mapstructure:"file_location"`
		FileMaxAge   int       `json:"file_max_age" mapstructure:"file_max_age"`
		Level        int8      `json:"level" mapstructure:"level"`
		SplitFile    SplitFile `json:"split_file" mapstructure:"split_file"`
	}

	SplitFile struct {
		Debug   bool `json:"debug" mapstructure:"debug"`
		Info    bool `json:"info" mapstructure:"info"`
		Warning bool `json:"warning" mapstructure:"warning"`
		Error   bool `json:"error" mapstructure:"error"`
		Fatal   bool `json:"fatal" mapstructure:"fatal"`
		Tdr     bool `json:"tdr" mapstructure:"tdr"`
	}
)

type (
	fileLogger struct {
		defaultLogger *zapLogger
		debugLogger   *zapLogger
		infoLogger    *zapLogger
		warningLogger *zapLogger
		errorLogger   *zapLogger
		fatalLogger   *zapLogger
		tdrLogger     *zapLogger
	}
)

const (
	Info  Level = "INFO"
	Debug Level = "DEBUG"
	Error Level = "ERROR"

	JSONFormatter Formatter = "JSON"
)

const (
	externalIDKey     = "external_id"
	journeyIDKey      = "journey_id"
	chainIDKey        = "chain_id"
	responseTimeKey   = "rt"
	responseURI       = "uri"
	headerKey         = "header"
	requestKey        = "req"
	responseKey       = "resp"
	errorKey          = "error"
	additionalDataKey = "additional_data"
	levelKey          = "level"
)

func NewLogger(config *LogOption) Logger {
	fmt.Println("Try NewLogger File...")

	if config == nil {
		panic("logger file config is nil")
	}

	logger := &fileLogger{
		defaultLogger: createLogger(config.Stdout, config.Level, config.FileLocation, config.FileMaxAge),
	}

	if !config.Stdout && config.SplitFile.Debug {
		logger.debugLogger = createLogger(config.Stdout, config.Level, config.FileLocation+".debug", config.FileMaxAge)
	}
	if !config.Stdout && config.SplitFile.Info {
		logger.infoLogger = createLogger(config.Stdout, config.Level, config.FileLocation+".info", config.FileMaxAge)
	}
	if !config.Stdout && config.SplitFile.Warning {
		logger.warningLogger = createLogger(config.Stdout, config.Level, config.FileLocation+".warning", config.FileMaxAge)
	}
	if !config.Stdout && config.SplitFile.Error {
		logger.errorLogger = createLogger(config.Stdout, config.Level, config.FileLocation+".error", config.FileMaxAge)
	}
	if !config.Stdout && config.SplitFile.Fatal {
		logger.fatalLogger = createLogger(config.Stdout, config.Level, config.FileLocation+".fatal", config.FileMaxAge)
	}
	if !config.Stdout && config.SplitFile.Tdr {
		logger.tdrLogger = createLogger(config.Stdout, config.Level, config.FileLocation+".tdr", config.FileMaxAge)
	}

	if logger.debugLogger == nil {
		logger.debugLogger = logger.defaultLogger
	}
	if logger.infoLogger == nil {
		logger.infoLogger = logger.defaultLogger
	}
	if logger.warningLogger == nil {
		logger.warningLogger = logger.defaultLogger
	}
	if logger.errorLogger == nil {
		logger.errorLogger = logger.defaultLogger
	}
	if logger.fatalLogger == nil {
		logger.fatalLogger = logger.defaultLogger
	}
	if logger.tdrLogger == nil {
		logger.tdrLogger = logger.defaultLogger
	}

	return logger
}

func (c *fileLogger) Info(args ...interface{}) {
	c.InfoWithCtx(context.Background(), fmt.Sprint(args...))
}

func (c *fileLogger) Infof(format string, args ...interface{}) {
	c.InfofWithCtx(context.Background(), format, args...)
}

func (c *fileLogger) Debug(args ...interface{}) {
	c.DebugWithCtx(context.Background(), fmt.Sprint(args...))
}

func (c *fileLogger) Debugf(format string, args ...interface{}) {
	c.DebugfWithCtx(context.Background(), format, args...)
}

func (c *fileLogger) Error(args ...interface{}) {
	c.ErrorWithCtx(context.Background(), fmt.Sprint(args...))
}

func (c *fileLogger) Errorf(format string, args ...interface{}) {
	c.ErrorfWithCtx(context.Background(), format, args...)
}

func (c *fileLogger) Warning(args ...interface{}) {
	c.WarnWithCtx(context.Background(), fmt.Sprint(args...))
}

func (c *fileLogger) Warningf(format string, args ...interface{}) {
	c.WarnfWithCtx(context.Background(), format, args...)
}

func (c *fileLogger) Fatal(args ...interface{}) {
	c.FatalWithCtx(context.Background(), fmt.Sprint(args...))
}

func (c *fileLogger) Fatalf(format string, args ...interface{}) {
	c.FatalfWithCtx(context.Background(), format, args...)
}

func (c *fileLogger) Print(args ...interface{}) {
	c.infoLogger.Print(args...)
}

func (c *fileLogger) Printf(format string, args ...interface{}) {
	c.infoLogger.Printf(format, args...)
}

func (c *fileLogger) Println(args ...interface{}) {
	c.infoLogger.Println(args...)
}

func (c *fileLogger) Instance() interface{} {
	return c.defaultLogger
}

func (c *fileLogger) DebugWithCtx(ctx context.Context, message string, fields ...Field) {
	c.debugLogger.Debug(ctx, message, fields...)
}

func (c *fileLogger) DebugfWithCtx(ctx context.Context, format string, args ...interface{}) {
	c.debugLogger.Debugf(ctx, format, args...)
}

func (c *fileLogger) InfoWithCtx(ctx context.Context, message string, fields ...Field) {
	c.infoLogger.Info(ctx, message, fields...)
}

func (c *fileLogger) InfofWithCtx(ctx context.Context, format string, args ...interface{}) {
	c.infoLogger.Infof(ctx, format, args...)
}

func (c *fileLogger) WarnWithCtx(ctx context.Context, message string, fields ...Field) {
	c.warningLogger.Warn(ctx, message, fields...)
}

func (c *fileLogger) WarnfWithCtx(ctx context.Context, format string, args ...interface{}) {
	c.warningLogger.Warnf(ctx, format, args...)
}

func (c *fileLogger) ErrorWithCtx(ctx context.Context, message string, fields ...Field) {
	c.errorLogger.Error(ctx, message, fields...)
}

func (c *fileLogger) ErrorfWithCtx(ctx context.Context, format string, args ...interface{}) {
	c.errorLogger.Errorf(ctx, format, args...)
}

func (c *fileLogger) FatalWithCtx(ctx context.Context, message string, fields ...Field) {
	c.fatalLogger.Fatal(ctx, message, fields...)
}

func (c *fileLogger) FatalfWithCtx(ctx context.Context, format string, args ...interface{}) {
	c.fatalLogger.Fatalf(ctx, format, args...)
}

func (c *fileLogger) Summary(tdr LogSummary) {
	c.tdrLogger.Summary(tdr)
}
