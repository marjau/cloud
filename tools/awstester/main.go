package main

import (
	"log"
)

func main() {
	log.Println("Starting testing AWS services...")

	cfg, err := newAWSConfig()
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// TODO: retrieve specific testers from environment variables
	testers := []AWSTester{
		newS3Tester(cfg),
		newSQSTester(cfg),
		newSNSTester(cfg),
		newDynamoDBTester(cfg),
	}

	// TODO: apply concurrency
	errs := make([]error, 0)
	for _, tt := range testers {
		if err := tt.RunTests(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		for i, e := range errs {
			log.Printf("  %d. Failed running test: %s", i+1, e)
		}
		log.Fatalln("Errors occurred when running the AWS services tests.")
	}
	log.Println("All tests passed")
}
