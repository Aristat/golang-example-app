package logger

import (
	"context"

	"go.uber.org/zap"
)

type (
	Fields map[string]interface{}
	Tags   []string
)

type opts struct {
	args   []interface{}
	tags   Tags
	fields Fields
}

// Option is func hook for underling logic call
type Option func(*opts) error

// Zap is uber/zap logger implemented of Logger interface
type Zap struct {
	ctx    context.Context
	cfg    Config
	logger *zap.SugaredLogger
	fields map[string]interface{}
	tags   []string
}

// Printf is like fmt.Printf, push to log entry with debug level
func (z *Zap) Printf(format string, a ...interface{}) {
	z.Debug(format, Args(a...))
}

// Verbose should return true when verbose logging output is wanted
func (z *Zap) Verbose() bool {
	return z.cfg.Verbose
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
		tags  []string
	)
	// tags
	if len(z.tags) > 0 {
		tags = make([]string, len(z.tags))
		copy(tags, z.tags)
	}
	if len(opts.tags) > 0 {
		tags = append(tags, opts.tags...)
	}
	if len(tags) > 0 {
		wargs = append(wargs, "tags", tags)
	}
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
func (z *Zap) WithFields(fields Fields) *Zap {
	nz := &Zap{}
	copyZap(nz, z, nil, fields)
	return nz
}

// WithTags create new instance with tags
func (z *Zap) WithTags(tags Tags) *Zap {
	nz := &Zap{}
	copyZap(nz, z, tags, nil)
	return nz
}

// Args returns func hook a logger for replace fmt placeholders on represent values
func Args(a ...interface{}) Option {
	return func(f *opts) error {
		f.args = a
		return nil
	}
}

func copyZap(dst, src *Zap, tags []string, fields map[string]interface{}) {
	// tags
	var cTags []string
	cTags = make([]string, len(src.tags))
	copy(cTags, src.tags)
	dst.tags = cTags
	if tags != nil {
		dst.tags = append(dst.tags, tags...)
	}
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
		cfg.EncoderConfig.LevelKey = ""
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}
		logger, _ = cfg.Build(zap.AddCallerSkip(2), zap.AddStacktrace(zap.PanicLevel))
	} else {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.LevelKey = ""
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}
		logger, _ = cfg.Build(zap.AddCallerSkip(2))
	}
	go func(logger *zap.Logger) {
		<-ctx.Done()
		_ = logger.Sync()
	}(logger)
	return &Zap{ctx: ctx, cfg: cfg, logger: logger.Sugar()}
}
