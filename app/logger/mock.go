package logger

import (
	"context"
)

// Entity represent log struct with all assets
type Entity struct {
	Level  Level
	Fields Fields
	Args   []interface{}
	Format string
}

// Mock is the logger with stubbed methods
type Mock struct {
	ctx     context.Context
	ch      chan Entity
	discard bool
	cfg     Config
	fields  map[string]interface{}
}

// Printf is like fmt.Printf, push to log entry with debug level
func (m *Mock) Printf(format string, a ...interface{}) {
	m.Debug(format, Args(a...))
}

// Emergency push to log entry with emergency level & throw panic
func (m *Mock) Emergency(format string, opts ...Option) {
	m.Log(LevelEmergency, format, opts...)
}

// Alert push to log entry with alert level
func (m *Mock) Alert(format string, opts ...Option) {
	m.Log(LevelAlert, format, opts...)
}

// Critical push to log entry with critical level
func (m *Mock) Critical(format string, opts ...Option) {
	m.Log(LevelCritical, format, opts...)
}

// Error push to log entry with error level
func (m *Mock) Error(format string, opts ...Option) {
	m.Log(LevelError, format, opts...)
}

// Warning push to log entry with warning level
func (m *Mock) Warning(format string, opts ...Option) {
	m.Log(LevelWarning, format, opts...)
}

// Notice push to log entry with notice level
func (m *Mock) Notice(format string, opts ...Option) {
	m.Log(LevelNotice, format, opts...)
}

// Info push to log entry with info level
func (m *Mock) Info(format string, opts ...Option) {
	m.Log(LevelInfo, format, opts...)
}

// Debug push to log entry with debug level
func (m *Mock) Debug(format string, opts ...Option) {
	m.Log(LevelDebug, format, opts...)
}

// Write push to log entry with debug level
func (m *Mock) Write(p []byte) (n int, err error) {
	m.Debug(string(p))
	return len(p), nil
}

// Log push to log with specified level
func (m *Mock) Log(level Level, format string, o ...Option) {
	if !m.discard {
		return
	}
	opts := &opts{}
	for _, option := range o {
		_ = option(opts)
	}
	var (
		fields = map[string]interface{}{}
	)
	// fields
	for k, v := range m.fields {
		fields[k] = v
	}
	for k, v := range opts.fields {
		fields[k] = v
	}

	go func() {
		m.ch <- Entity{
			Level:  level,
			Format: format,
			Args:   opts.args,
			Fields: fields,
		}
	}()
}

// WithFields create new instance with fields
func (m *Mock) WithFields(fields Fields) Logger {
	nm := &Mock{}
	copyMock(nm, m, fields)
	return nm
}

func copyMock(dst, src *Mock, fields map[string]interface{}) {
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
	dst.discard = src.discard
	dst.ch = src.ch
}

// Catch returns channel of entity structure for testing event content
func (m *Mock) Catch() <-chan Entity {
	return m.ch
}

// NewMock returns mock instance implemented of Logger interface
func newMock(ctx context.Context, cfg Config, discard bool) *Mock {
	return &Mock{
		ctx:     ctx,
		cfg:     cfg,
		ch:      make(chan Entity),
		discard: discard,
	}
}
