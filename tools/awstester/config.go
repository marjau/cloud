package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	AWS_PROFILE_DEFAULT      = "default"
	AWS_REGION_DEFAULT       = "us-east-1"
	AWS_ENDPOINT_URL_DEFAULT = "http://localhost:4566"
)

// newAWSConfig loads the AWS configuration for SDK v2
func newAWSConfig() (aws.Config, error) {
	region := getEnv("AWS_REGION", AWS_REGION_DEFAULT)
	endpoint := getEnv("AWS_ENDPOINT_URL", AWS_ENDPOINT_URL_DEFAULT)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint, SigningRegion: region}, nil
			}),
		),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	return cfg, nil
}

// getEnv is a helper function to get environment variables with a fallback value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
