package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

type AWSTester interface {
	RunTests() error
}

type awsTester struct {
	cfg    aws.Config
	logger *prefixedLogger
	// t *testing.T
}
