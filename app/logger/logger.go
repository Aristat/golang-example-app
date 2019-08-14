package logger

type (
	Fields map[string]interface{}
)

// Level represent RFC5424 logger severity
type Level int8

// String implements interface Stringer
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelNotice:
		return "notice"
	case LevelWarning:
		return "warning"
	case LevelError:
		return "error"
	case LevelCritical:
		return "critical"
	case LevelAlert:
		return "alert"
	case LevelEmergency:
		return "emergency"
	}
	return "unknown"
}

// FromString set logger level from string representation
func (l *Level) FromString(level string) Level {
	switch level {
	case "debug":
		*l = LevelDebug
	case "info":
		*l = LevelInfo
	case "notice":
		*l = LevelNotice
	case "warning":
		*l = LevelWarning
	case "error":
		*l = LevelError
	case "critical":
		*l = LevelCritical
	case "alert":
		*l = LevelAlert
	case "emergency":
		*l = LevelEmergency
	}
	return *l
}

const (
	// Debug: debug-level messages
	LevelDebug Level = 7
	// Informational: informational messages
	LevelInfo Level = 6
	// Notice: normal but significant condition
	LevelNotice Level = 5
	// Warning: warning conditions
	LevelWarning Level = 4
	// Error: error conditions
	LevelError Level = 3
	// Critical: critical conditions
	LevelCritical Level = 2
	// Alert: action must be taken immediately
	LevelAlert Level = 1
	// Emergency: system is unusable
	LevelEmergency Level = 0
)

type opts struct {
	args   []interface{}
	fields Fields
}

// Option is func hook for underling logic call
type Option func(*opts) error

// Args returns func hook a logger for replace fmt placeholders on represent values
func Args(a ...interface{}) Option {
	return func(f *opts) error {
		f.args = a
		return nil
	}
}

// WithFields returns func hook a logger for adding fields for call
func WithFields(fields Fields) Option {
	return func(f *opts) error {
		f.fields = fields
		return nil
	}
}

// Config is a general logger config settings
type Config struct {
	Debug bool
}

// Logger is the interface for logger client
type Logger interface {
	// Printf is like fmt.Printf, push to log entry with debug level
	Printf(format string, a ...interface{})
	// Emergency push to log entry with emergency level & throw panic
	Emergency(format string, opts ...Option)
	// Alert push to log entry with alert level
	Alert(format string, opts ...Option)
	// Critical push to log entry with critical level
	Critical(format string, opts ...Option)
	// Error push to log entry with error level
	Error(format string, opts ...Option)
	// Warning push to log entry with warning level
	Warning(format string, opts ...Option)
	// Notice push to log entry with notice level
	Notice(format string, opts ...Option)
	// Info push to log entry with info level
	Info(format string, opts ...Option)
	// Debug push to log entry with debug level
	Debug(format string, opts ...Option)
	// Write push to log entry with debug level
	Write(p []byte) (n int, err error)
	// Log push to log with specified level
	Log(level Level, format string, opts ...Option)
	// WithFields create new instance with fields
	WithFields(fields Fields) Logger
}
