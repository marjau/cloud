package logger

import (
	"fmt"
	"log"
	"os"
)

// PrefixedLogger wraps the standard log.Logger
// and includes an additional value in each log message.
type PrefixedLogger struct {
	logger *log.Logger
	prefix string
}

// NewPrefixedLogger initializes a new prefixedLogger with a specific prefix.
func NewPrefixedLogger(prefix string) *PrefixedLogger {
	return &PrefixedLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		prefix: prefix,
	}
}

// Log logs a message with the custom prefix.
func (pl *PrefixedLogger) Log(message string) {
	pl.logger.Printf("[%v] %v", pl.prefix, message)
}

// logf logs a formatted message with the custom prefix.
func (pl *PrefixedLogger) Logf(format string, args ...interface{}) {
	pl.logger.Printf(fmt.Sprintf("[%v] %v", pl.prefix, format), args...)
}

func (pl *PrefixedLogger) Fataln(message string) {
	pl.logger.Printf("[%v][FATAL] %v", pl.prefix, message)
	os.Exit(1)
}

func (pl *PrefixedLogger) Fatalf(format string, args ...interface{}) {
	pl.logger.Printf(fmt.Sprintf("[%v][FATAL] %v", pl.prefix, format), args...)
	os.Exit(1)
}
