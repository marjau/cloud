package config

import (
	"context"
	"errors"
	"myapps/awstester/internal/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	AWS_PROFILE_DEFAULT      = "default"
	AWS_REGION_DEFAULT       = "us-east-1"
	AWS_ENDPOINT_URL_DEFAULT = "http://localhost:4566"
)

var (
	errRegionEmpty  = errors.New("AWS Region is empty")
	errEnpointEmpty = errors.New("AWS EndPoint is empty")
)

// NewAWSConfig loads the AWS configuration from the given region
func NewAWSConfig(region string) (aws.Config, error) {
	l := logger.NewPrefixedLogger("CONFIG")
	l.Log("Loading AWS configurations")

	if region == "" {
		return aws.Config{}, errRegionEmpty
	}
	l.Logf("  Region: %v", region)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}

// NewAWSConfigWithEnpoint loads the AWS configuration
func NewAWSConfigWithEnpoint(region string, endpoint string) (aws.Config, error) {
	l := logger.NewPrefixedLogger("CONFIG")
	l.Log("Loading AWS configurations with endpoint")

	if region == "" {
		return aws.Config{}, errRegionEmpty
	}
	if endpoint == "" {
		return aws.Config{}, errEnpointEmpty
	}
	l.Logf("  Region: %v", region)
	l.Logf("  Endpoint: %v", endpoint)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint, SigningRegion: region}, nil
			}),
		),
	)
	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}
