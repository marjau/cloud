package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// TESTER *****************

type Tester interface {
	RunTests() error
	Clean()
}

type tester struct {
	config AWSConfig
	logger *prefixedLogger
}

// SESSION *****************

func newAWSSession(cfg AWSConfig) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.Region),
		Endpoint:         aws.String(cfg.EndpointURL),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	// log.Println("Session created sucessfully")
	return sess, nil
}

// LOGGER *****************

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
	pl.logger.Println(pl.prefix + ": " + message)
}

// logf logs a formatted message with the custom prefix.
func (pl *prefixedLogger) logf(format string, args ...interface{}) {
	pl.logger.Printf(fmt.Sprintf("%v: %v", pl.prefix, format), args...)
}
