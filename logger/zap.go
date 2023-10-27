package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	ZapOption func(*zapLogger) error

	zapLogger struct {
		writers []io.Writer
		closer  []io.Closer
		zapLog  *zap.Logger
		level   int8
	}
)

func WithStdout() ZapOption {
	return func(logger *zapLogger) error {
		// Wire STD output for both type
		logger.writers = append(logger.writers, os.Stdout)
		return nil
	}
}

func WithFileOutput(location string, maxAge int) ZapOption {
	return func(logger *zapLogger) error {
		output, err := rotateLogs.New(
			location+".%Y%m%d",
			rotateLogs.WithLinkName(location),
			rotateLogs.WithMaxAge(time.Duration(maxAge)*24*time.Hour),
			rotateLogs.WithRotationTime(time.Hour),
		)

		if err != nil {
			return fmt.Errorf("sys file error: %w", err)
		}

		// Wire SYS config only in sys
		logger.writers = append(logger.writers, output)
		logger.closer = append(logger.closer, output)
		return nil
	}
}

// WithCustomWriter add custom writer, so you can write using any storage method
// without waiting this package to be updated.
func WithCustomWriter(writer io.WriteCloser) ZapOption {
	return func(logger *zapLogger) error {
		if writer == nil {
			return fmt.Errorf("writer is nil")
		}

		// wire custom writer to log
		logger.writers = append(logger.writers, writer)
		logger.closer = append(logger.closer, writer)
		return nil
	}
}

// WithLevel set level of logger
func WithLevel(level int8) ZapOption {
	return func(logger *zapLogger) error {
		logger.level = level
		return nil
	}
}

func NewZapLogger(opts ...ZapOption) (*zapLogger, error) {
	defaultLogger := &zapLogger{
		writers: make([]io.Writer, 0),
	}

	for _, o := range opts {
		if err := o(defaultLogger); err != nil {
			return nil, err
		}
	}

	// use stdout only when writer is not specified
	if len(defaultLogger.writers) <= 0 {
		defaultLogger.writers = append(defaultLogger.writers, zapcore.AddSync(os.Stdout))
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "time",
		MessageKey: "message",
		LineEnding: zapcore.DefaultLineEnding,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(d.Nanoseconds() / 1000000)
		},
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.999"))
		},
	}

	encoding := zapcore.NewJSONEncoder(encoderConfig)
	// set logger here instead in options to make easy and consistent initiation
	// set multiple writer as already set in options
	zapWriters := make([]zapcore.WriteSyncer, 0)
	for _, writer := range defaultLogger.writers {
		if writer == nil {
			continue
		}

		zapWriters = append(zapWriters, zapcore.AddSync(writer))
	}

	core := zapcore.NewCore(encoding, zapcore.NewMultiWriteSyncer(zapWriters...), zapcore.Level(defaultLogger.level))
	defaultLogger.zapLog = zap.New(core)

	return defaultLogger, nil
}

func (d *zapLogger) Close() error {
	if d.closer == nil {
		return nil
	}

	var err error
	for _, closer := range d.closer {
		if closer == nil {
			continue
		}

		if e := closer.Close(); e != nil {
			err = fmt.Errorf("%w: %q", e, err)
		}
	}

	return err
}

func (d *zapLogger) Debug(ctx context.Context, message string, fields ...Field) {
	zapLogs := []zap.Field{
		zap.String(levelKey, "debug"),
	}
	zapLogs = append(zapLogs, ctxToLog(ctx)...)
	zapLogs = append(zapLogs, appendToLog(fields...)...)
	d.zapLog.Debug(message, zapLogs...)
}

func (d *zapLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	zapLogs := []zap.Field{
		zap.String(levelKey, "debug"),
	}
	zapLogs = append(zapLogs, ctxToLog(ctx)...)
	d.zapLog.Debug(fmt.Sprintf(format, args...), zapLogs...)
}

func (d *zapLogger) Info(ctx context.Context, message string, fields ...Field) {
	zapLogs := []zap.Field{
		zap.String(levelKey, "info"),
	}
	zapLogs = append(zapLogs, ctxToLog(ctx)...)
	zapLogs = append(zapLogs, appendToLog(fields...)...)
	d.zapLog.Info(message, zapLogs...)
}

func (d *zapLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	zapLogs := []zap.Field{
		zap.String(levelKey, "info"),
	}
	zapLogs = append(zapLogs, ctxToLog(ctx)...)
	d.zapLog.Info(fmt.Sprintf(format, args...), zapLogs...)
}

func (d *zapLogger) Warn(ctx context.Context, message string, fields ...Field) {
	zapLogs := []zap.Field{
		zap.String(levelKey, "warn"),
	}
	zapLogs = append(zapLogs, ctxToLog(ctx)...)
	zapLogs = append(zapLogs, appendToLog(fields...)...)
	d.zapLog.Warn(message, zapLogs...)
}

func (d *zapLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	zapLogs := []zap.Field{
		zap.String(levelKey, "warn"),
	}
	zapLogs = append(zapLogs, ctxToLog(ctx)...)
	d.zapLog.Warn(fmt.Sprintf(format, args...), zapLogs...)
}

func (d *zapLogger) Error(ctx context.Context, message string, fields ...Field) {
	zapLogs := []zap.Field{
		zap.String(levelKey, "error"),
	}
	zapLogs = append(zapLogs, ctxToLog(ctx)...)
	zapLogs = append(zapLogs, appendToLog(fields...)...)
	d.zapLog.Error(message, zapLogs...)
}

func (d *zapLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	zapLogs := []zap.Field{
		zap.String(levelKey, "error"),
	}
	zapLogs = append(zapLogs, ctxToLog(ctx)...)
	d.zapLog.Error(fmt.Sprintf(format, args...), zapLogs...)
}

func (d *zapLogger) Fatal(ctx context.Context, message string, fields ...Field) {
	zapLogs := []zap.Field{
		zap.String(levelKey, "fatal"),
	}
	zapLogs = append(zapLogs, ctxToLog(ctx)...)
	zapLogs = append(zapLogs, appendToLog(fields...)...)
	d.zapLog.Fatal(message, zapLogs...)
}

func (d *zapLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	zapLogs := []zap.Field{
		zap.String(levelKey, "fatal"),
	}
	zapLogs = append(zapLogs, ctxToLog(ctx)...)
	d.zapLog.Fatal(fmt.Sprintf(format, args...), zapLogs...)
}

func (d *zapLogger) Print(args ...interface{}) {
	d.zapLog.Info(fmt.Sprint(args...))
}

func (d *zapLogger) Printf(format string, args ...interface{}) {
	d.zapLog.Info(fmt.Sprintf(format, args...))
}

func (d *zapLogger) Println(args ...interface{}) {
	d.zapLog.Info(fmt.Sprint(args...))
}

func (d *zapLogger) Summary(tdr LogSummary) {
	fields := []zap.Field{
		zap.String(levelKey, "tdr"),
		zap.String(externalIDKey, tdr.ExternalID),
		zap.String(journeyIDKey, tdr.JourneyID),
		zap.String(chainIDKey, tdr.ChainID),
		zap.Int64(responseTimeKey, tdr.RespTime),
		zap.String(responseURI, tdr.URI),
		formatLog(headerKey, tdr.Header),
		formatLog(requestKey, tdr.Request),
		formatLog(responseKey, tdr.Response),
		zap.String(errorKey, tdr.Error),
		formatLog(additionalDataKey, tdr.AdditionalData),
	}

	d.zapLog.Info("Summary", fields...)
}
