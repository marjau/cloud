package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"myapps/awstester/internal/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const s3SvcName = "S3"

type S3Tester struct {
	awsTester
	svc *s3.Client
}

func NewS3Tester(cfg aws.Config) *S3Tester {
	return &S3Tester{
		awsTester: awsTester{
			cfg:    cfg,
			logger: logger.NewPrefixedLogger(s3SvcName),
		},
		svc: s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true // TODO: move to docker hosts approach
		}),
	}
}

func (st S3Tester) GetName() string {
	return s3SvcName
}

func (st S3Tester) Run() error {
	st.logger.Log("Start Testing...")
	bucketName := "test-bucket"
	objectKey := "test-file.txt"
	fileContent := "Hello, this is a test file!"
	defer st.clean(bucketName)

	if err := st.createBucket(bucketName); err != nil {
		return err
	}
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

	st.logger.Log("Testing Completed.")
	return nil
}

func (st S3Tester) clean(bucketName string) {
	st.logger.Log("Start Cleaning Testing...")
	if err := st.DeleteBucket(bucketName); err != nil {
		var noSuchBucketErr *types.NoSuchBucket

		if errors.As(err, &noSuchBucketErr) {
			st.logger.Logf("INFO: Bucket %v not found", bucketName)
		} else {
			st.logger.Logf("ERROR: %v", err)
		}
	}
	st.logger.Log("Cleaning Testing Complete")
}

func (st S3Tester) createBucket(bucketName string) error {
	st.logger.Logf("Creating Bucket: %v", bucketName)

	_, err := st.svc.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket %v: %w", bucketName, err)
	}

	st.logger.Log(fmt.Sprintf("Bucket %v created successfully", bucketName))
	return nil
}

func (st S3Tester) listBuckets() error {
	st.logger.Log("Listing Buckets")

	result, err := st.svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("failed to list buckets: %w", err)
	}

	st.logger.Log("Buckets:")
	for i, bucket := range result.Buckets {
		st.logger.Logf("  %d. %s", i+1, aws.ToString(bucket.Name))
	}
	return nil
}

func (st S3Tester) UploadFile(bucketName, objectKey, content string) error {
	st.logger.Logf("Uploading File to Bucket: %v, Key: %v", bucketName, objectKey)

	_, err := st.svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader([]byte(content)),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file %v to bucket %v: %w", objectKey, bucketName, err)
	}

	st.logger.Log("File uploaded successfully")
	return nil
}

func (st S3Tester) DownloadFile(bucketName, objectKey string) error {
	st.logger.Logf("Downloading File from Bucket: %v, Key: %v", bucketName, objectKey)

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

	st.logger.Logf("Downloaded file content: %s", string(body))
	return nil
}

func (st S3Tester) ListObjects(bucketName string) error {
	st.logger.Logf("Listing Objects in Bucket: %v", bucketName)

	result, err := st.svc.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to list objects in bucket %v: %w", bucketName, err)
	}

	st.logger.Log("Objects:")
	for i, item := range result.Contents {
		st.logger.Logf("\t%d. %s", i+1, aws.ToString(item.Key))
	}
	return nil
}

func (st S3Tester) DeleteBucket(bucketName string) error {
	st.logger.Logf("Deleting Bucket: %v", bucketName)

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

	st.logger.Log("Bucket deleted successfully")
	return nil
}

func (st S3Tester) deleteAllObjects(bucketName string) error {
	st.logger.Logf("Deleting all objects in Bucket: %v", bucketName)

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
		st.logger.Logf("Deleted object %v", aws.ToString(item.Key))
	}
	return nil
}
