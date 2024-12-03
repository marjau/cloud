package service

import (
	"myapps/awstester/internal/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type AWSTester interface {
	GetName() string
	Run() error
}

type awsTester struct {
	cfg    aws.Config
	logger *logger.PrefixedLogger
	// t *testing.T
}
