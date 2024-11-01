package main

import (
	"fmt"
	"log"
	"os"
)

// prefixedLogger wraps the standard log.Logger
// and includes an additional value in each log message.
type prefixedLogger struct {
	logger *log.Logger
	prefix string
}

// newPrefixedLogger initializes a new prefixedLogger with a specific prefix.
func newPrefixedLogger(prefix string) *prefixedLogger {
	return &prefixedLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		prefix: prefix,
	}
}

// log logs a message with the custom prefix.
func (pl *prefixedLogger) log(message string) {
	pl.logger.Printf("[%v] %v", pl.prefix, message)
}

// logf logs a formatted message with the custom prefix.
func (pl *prefixedLogger) logf(format string, args ...interface{}) {
	pl.logger.Printf(fmt.Sprintf("[%v] %v", pl.prefix, format), args...)
}
