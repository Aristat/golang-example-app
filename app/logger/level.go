package logger

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
