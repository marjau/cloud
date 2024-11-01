package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const s3SvcName = "S3"

type s3Tester struct {
	awsTester
	svc *s3.Client
}

func newS3Tester(cfg aws.Config) *s3Tester {
	return &s3Tester{
		awsTester: awsTester{
			cfg:    cfg,
			logger: newPrefixedLogger(s3SvcName),
		},
		svc: s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true // TODO: move to docker hosts approach
		}),
	}
}

func (st *s3Tester) RunTests() error {
	st.logger.log("Start Testing...")
	bucketName := "test-bucket"
	objectKey := "test-file.txt"
	fileContent := "Hello, this is a test file!"

	if err := st.createBucket(bucketName); err != nil {
		return err
	}
	defer st.clean(bucketName)

	if err := st.listBuckets(); err != nil {
		return err
	}
	if err := st.UploadFile(bucketName, objectKey, fileContent); err != nil {
		return err
	}
	if err := st.ListObjects(bucketName); err != nil {
		return err
	}
	if err := st.DownloadFile(bucketName, objectKey); err != nil {
		return err
	}
	if err := st.DeleteBucket(bucketName); err != nil {
		return err
	}

	st.logger.log("Testing Completed.")
	return nil
}

func (st *s3Tester) clean(bucketName string) {
	st.logger.log("Start Cleaning Testing...")
	if err := st.DeleteBucket(bucketName); err != nil {
		var noSuchBucketErr *types.NoSuchBucket

		if errors.As(err, &noSuchBucketErr) {
			st.logger.logf("INFO: Bucket %v not found", bucketName)
		} else {
			st.logger.logf("ERROR: %v", err)
		}
	}
	st.logger.log("Cleaning Testing Complete")
}

func (st *s3Tester) createBucket(bucketName string) error {
	st.logger.logf("Creating Bucket: %v", bucketName)

	_, err := st.svc.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket %v: %w", bucketName, err)
	}

	st.logger.log(fmt.Sprintf("Bucket %v created successfully", bucketName))
	return nil
}

func (st *s3Tester) listBuckets() error {
	st.logger.log("Listing Buckets")

	result, err := st.svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("failed to list buckets: %w", err)
	}

	st.logger.log("Buckets:")
	for i, bucket := range result.Buckets {
		st.logger.logf("  %d. %s", i+1, aws.ToString(bucket.Name))
	}
	return nil
}

func (st *s3Tester) UploadFile(bucketName, objectKey, content string) error {
	st.logger.logf("Uploading File to Bucket: %v, Key: %v", bucketName, objectKey)

	_, err := st.svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader([]byte(content)),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file %v to bucket %v: %w", objectKey, bucketName, err)
	}

	st.logger.log("File uploaded successfully")
	return nil
}

func (st *s3Tester) DownloadFile(bucketName, objectKey string) error {
	st.logger.logf("Downloading File from Bucket: %v, Key: %v", bucketName, objectKey)

	result, err := st.svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to download file %v from bucket %v: %w", objectKey, bucketName, err)
	}
	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	st.logger.logf("Downloaded file content: %s", string(body))
	return nil
}

func (st *s3Tester) ListObjects(bucketName string) error {
	st.logger.logf("Listing Objects in Bucket: %v", bucketName)

	result, err := st.svc.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to list objects in bucket %v: %w", bucketName, err)
	}

	st.logger.log("Objects:")
	for i, item := range result.Contents {
		st.logger.logf("\t%d. %s", i+1, aws.ToString(item.Key))
	}
	return nil
}

func (st *s3Tester) DeleteBucket(bucketName string) error {
	st.logger.logf("Deleting Bucket: %v", bucketName)

	// Delete all objects in the bucket before deleting the bucket
	if err := st.deleteAllObjects(bucketName); err != nil {
		return err
	}

	_, err := st.svc.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete bucket %v: %w", bucketName, err)
	}

	st.logger.log("Bucket deleted successfully")
	return nil
}

func (st *s3Tester) deleteAllObjects(bucketName string) error {
	st.logger.logf("Deleting all objects in Bucket: %v", bucketName)

	// List objects to delete
	result, err := st.svc.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to list objects in bucket %v: %w", bucketName, err)
	}

	for _, item := range result.Contents {
		_, err := st.svc.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    item.Key,
		})
		if err != nil {
			return fmt.Errorf("failed to delete object %v: %w", aws.ToString(item.Key), err)
		}
		st.logger.logf("Deleted object %v", aws.ToString(item.Key))
	}
	return nil
}
