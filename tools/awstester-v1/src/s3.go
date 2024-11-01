package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Tester struct {
	tester
	svc *s3.S3
}

func NewS3Tester(cfg AWSConfig, sess *session.Session) *S3Tester {
	return &S3Tester{
		tester: tester{
			config: cfg,
			logger: newPrefixedLogger("S3"),
		},
		svc: s3.New(sess),
	}
}

func (st S3Tester) RunTests() error {
	st.logger.log("Start Testing...")
	if err := st.createBucket(); err != nil {
		return err
	}
	if err := st.listBuckets(); err != nil {
		return err
	}
	st.logger.log("Tests completed")
	return nil
}

func (st S3Tester) Clean() {
	st.logger.log("Cleaning tests...")
	st.logger.log("Tests cleaned")
}

func (st S3Tester) listBuckets() error {
	st.logger.log("Testing List Buckets")

	result, err := st.svc.ListBuckets(nil)
	if err != nil {
		log.Fatalf("Failed to list buckets: %v", err)
	}

	st.logger.log("Buckets:")
	for i, bucket := range result.Buckets {
		st.logger.logf("  %v. %s", i+1, aws.StringValue(bucket.Name))
	}
	return nil
}

func (st S3Tester) createBucket() error {
	bName := "test-bucket"
	st.logger.logf("Testing Create Bucket %v", bName)

	createBucketInput := &s3.CreateBucketInput{
		Bucket: aws.String(bName),
	}
	_, err := st.svc.CreateBucket(createBucketInput)
	if err != nil {
		return err
	}

	st.logger.log("Bucket created successfully")
	return nil
}
