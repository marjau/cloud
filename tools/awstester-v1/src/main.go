package main

import (
	"log"
	"os"
)

func main() {
	log.Println("Starting tests against AWS Services")

	cfg := NewAWSConfig()
	sess, err := newAWSSession(cfg)
	if err != nil {
		log.Fatalf("SESSION: Failed to create session: %v", err)
	}

	testers := []Tester{
		NewS3Tester(cfg, sess),
		NewSQSTester(cfg, sess),
		NewSNSTester(cfg, sess),
		NewDynamoDBTester(cfg, sess),
	}

	for _, tester := range testers {
		defer tester.Clean()
		err := tester.RunTests()
		if err != nil {
			log.Fatalf("Failed running tests: %v", err)
			os.Exit(1)
		}
	}
	log.Println("Tests against AWS Services completed")
}
