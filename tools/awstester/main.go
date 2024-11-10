package main

import (
	"myapps/awstester/internal/config"
	"myapps/awstester/internal/logger"
	svc "myapps/awstester/internal/service"
	"myapps/awstester/internal/utils"
)

const appName = "AWS-Tester"

func main() {
	l := logger.NewPrefixedLogger(appName)
	l.Log("Starting testing AWS services...")

	cfg, err := config.NewAWSConfigWithEnpoint(
		utils.GetEnv("AWS_REGION", config.AWS_REGION_DEFAULT),
		utils.GetEnv("AWS_ENDPOINT_URL", config.AWS_ENDPOINT_URL_DEFAULT),
	)
	if err != nil {
		l.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// TODO: retrieve specific awsTesters from environment variables
	awsTesters := []svc.AWSTester{
		svc.NewS3Tester(cfg),
		svc.NewSQSTester(cfg),
		svc.NewSNSTester(cfg),
		svc.NewDynamoDBTester(cfg),
	}

	// TODO: apply concurrency
	errs := make([]error, 0)
	for _, tester := range awsTesters {
		if err := tester.Run(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		for i, e := range errs {
			l.Logf("  %d. Failed running test: %s", i+1, e)
		}
		l.Fatalf("Errors occurred when running the AWS services tests.")
	}
	l.Log("All tests passed")
}
