package main

import (
	"os"
)

const (
	AWS_PROFILE_DEFAULT      = "localstack"
	AWS_REGION_DEFAULT       = "us-east-1"
	AWS_ENDPOINT_URL_DEFAULT = "http://localstack-main:4566"
)

type AWSConfig struct {
	Profile     string
	Region      string
	EndpointURL string
}

func NewAWSConfig() AWSConfig {
	l := newPrefixedLogger("Config")
	l.log("Loading AWS configurations...")

	profile := getEnv("AWS_PROFILE", AWS_REGION_DEFAULT)
	region := getEnv("AWS_REGION", AWS_REGION_DEFAULT)
	endpoint := getEnv("AWS_ENDPOINT_URL", AWS_ENDPOINT_URL_DEFAULT)

	l.logf("AWS PROFILE: %v", profile)
	l.logf("AWS REGION: %v", region)
	l.logf("AWS ENDPOINT URL: %v", endpoint)

	return AWSConfig{
		Profile:     profile,
		Region:      region,
		EndpointURL: endpoint,
	}
}

// getEnv is a helper function to get environment variables with a fallback value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
