package db

import (
	"bytes"
	"fmt"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"
)

type loggerAdapter struct {
	level  logger.Level
	logger *logger.Zap
}

// Print write string to output
func (l *loggerAdapter) Print(v ...interface{}) {
	var str bytes.Buffer
	for _, s := range v {
		if str.Len() > 0 {
			str.WriteString(" ")
		}
		str.WriteString(fmt.Sprint(s))
	}
	l.logger.Log(l.level, "%v", logger.Args(str.String()))
}

// NewLoggerAdapter returns instance adapter for services/logger
func NewLoggerAdapter(logger *logger.Zap, level logger.Level) *loggerAdapter {
	return &loggerAdapter{level: level, logger: logger}
}
