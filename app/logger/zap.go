package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Zap is uber/zap logger implemented of Logger interface
type Zap struct {
	ctx    context.Context
	cfg    Config
	logger *zap.SugaredLogger
	fields map[string]interface{}
}

// Printf is like fmt.Printf, push to log entry with debug level
func (z *Zap) Printf(format string, a ...interface{}) {
	z.Debug(format, Args(a...))
}

// Emergency push to log entry with emergency level & throw panic
func (z *Zap) Emergency(format string, opts ...Option) {
	z.Log(LevelEmergency, format, opts...)
}

// Alert push to log entry with alert level
func (z *Zap) Alert(format string, opts ...Option) {
	z.Log(LevelAlert, format, opts...)
}

// Critical push to log entry with critical level
func (z *Zap) Critical(format string, opts ...Option) {
	z.Log(LevelCritical, format, opts...)
}

// Error push to log entry with error level
func (z *Zap) Error(format string, opts ...Option) {
	z.Log(LevelError, format, opts...)
}

// Warning push to log entry with warning level
func (z *Zap) Warning(format string, opts ...Option) {
	z.Log(LevelWarning, format, opts...)
}

// Notice push to log entry with notice level
func (z *Zap) Notice(format string, opts ...Option) {
	z.Log(LevelNotice, format, opts...)
}

// Info push to log entry with info level
func (z *Zap) Info(format string, opts ...Option) {
	z.Log(LevelInfo, format, opts...)
}

// Debug push to log entry with debug level
func (z *Zap) Debug(format string, opts ...Option) {
	z.Log(LevelDebug, format, opts...)
}

// Write push to log entry with debug level
func (z *Zap) Write(p []byte) (n int, err error) {
	z.Debug(string(p))
	return len(p), nil
}

// Log push to log with specified level
func (z *Zap) Log(level Level, format string, o ...Option) {
	opts := &opts{}
	for _, option := range o {
		_ = option(opts)
	}
	var (
		wargs = []interface{}{"level", level.String()}
	)

	// fields
	for k, v := range z.fields {
		wargs = append(wargs, k, v)
	}
	for k, v := range opts.fields {
		wargs = append(wargs, k, v)
	}

	var logger = z.logger
	if len(wargs) > 0 {
		logger = logger.With(wargs...)
	}

	if len(opts.args) == 0 {
		var fn func(args ...interface{})
		switch level {
		default:
			fn = logger.Debug
		case LevelInfo, LevelNotice:
			fn = logger.Info
		case LevelWarning:
			fn = logger.Warn
		case LevelError, LevelCritical, LevelAlert:
			fn = logger.Error
		case LevelEmergency:
			fn = logger.Panic
		}
		fn(format)
	} else {
		var fn func(format string, args ...interface{})
		switch level {
		default:
			fn = logger.Debugf
		case LevelInfo, LevelNotice:
			fn = logger.Infof
		case LevelWarning:
			fn = logger.Warnf
		case LevelError, LevelCritical, LevelAlert:
			fn = logger.Errorf
		case LevelEmergency:
			fn = logger.Panicf
		}
		fn(format, opts.args...)
	}
}

// WithFields create new instance with fields
func (z *Zap) WithFields(fields Fields) Logger {
	nz := &Zap{}
	copyZap(nz, z, fields)
	return nz
}

func copyZap(dst, src *Zap, fields map[string]interface{}) {
	var cFields = map[string]interface{}{}
	// fields
	for k, v := range src.fields {
		cFields[k] = v
	}
	dst.fields = cFields
	if fields != nil {
		for k, v := range fields {
			dst.fields[k] = v
		}
	}
	dst.logger = src.logger
}

// NewZap returns zap logger
func NewZap(ctx context.Context, cfg Config) *Zap {
	var logger *zap.Logger

	if !cfg.Debug {
		cfg := zap.NewProductionConfig()
		logger, _ = cfg.Build(zap.AddCallerSkip(2), zap.AddStacktrace(zap.PanicLevel))
	} else {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = cfg.Build(zap.AddCallerSkip(2))
	}

	go func(logger *zap.Logger) {
		<-ctx.Done()
		_ = logger.Sync()
	}(logger)

	return &Zap{ctx: ctx, cfg: cfg, logger: logger.Sugar()}
}
