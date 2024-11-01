# AWS Tester App
AWS Tester is a Go-based application designed to test AWS services locally using a LocalStack image across different environments, allowing for combinations of Go and AWS SDK V2 versions to find the most suitable match for your needs. It supports various AWS SDK V2 services, including S3, SQS, SNS, and DynamoDB.

This application is ideal for local development and testing without incurring costs for actual AWS resources. By default, it is based on Go `1.15`, AWS SDK V2 [v1.21.2](https://github.com/aws/aws-sdk-go-v2/tree/v1.21.2), and Alpine 3.13 versions, with the aim of migrating from SDK V1 to V2.

## Features
- **Simulates AWS service**s locally with LocalStack.
- Supports basic operations for **S3** (buckets and objects), **SQS** (queues), **SNS** (topics), and **DynamoDB** (CRUD on tables).
- Easily configurable for multiple versions of Go and AWS SDK V2.

### Directory Structure
```
awstester/
├── aws_credentials/        # AWS CLI operations
├── versions/               # Directory for each Go and AWS SDK version configuration
├── main.go                 # Entry point of the application
├── config.go               # AWS configuration and environment setup
├── tester.go               # AWS Tester
├── dynamodb.go             # DynamoDB testing functions
├── s3.go                   # S3 testing functions
├── sns.go                  # SNS testing functions
├── sqs.go                  # SQS testing functions
├── logger                  # Custom logger
├── go.mod                  # Go module file specifying dependencies
└── go.sum                  # Dependency checksums
```

## Getting Started

### Prerequisites
Ensure **Docker** and **Docker** Compose both are installed and running to spin up LocalStack.

### Set Up Docker and LocalStack
Use the docker-compose.yml file to start LocalStack, which will simulate the required AWS services.

Start the LocalStack container:
```
docker-compose up -d localstack
```
Ensure that the localstack-main container is running and fully operational.
```
docker ps -a | grep localstack-main
```
```
docker logs -f localstack-main
```

### Run the Application
Configure the container with the specific GO version and OS to run the tests
```
docker-compose up -d go1.15-alpine
```

Enter the created work environment 
```
docker exec -it awstester-1.15alp ash
```

Replace the `go.mod` and `go.sum` files with the compatible Go and AWS SDK V2 versions from the `versions` folder:
```
cp -fv versions/v1.5_sdk1.21.2/* .
go mod download
```
You can test multiple versions of the AWS SDK V2 on the selected Go version for the work environment by directly modifying the `go.mod` file.

Compile and run the tests:
```
go run *.go
```

### Cleaning Up
```
docker-compose down
```

## Testing Individual Services
You can run tests for individual services by calling specific functions in the code. 

### DynamoDB:
- **CreateTable**: Creates a DynamoDB table with ID as the primary key.
- **ListTables:** Lists all DynamoDB tables in the account.
- **PutItem:** Adds an item to a specified table.
- **GetItem:** Retrieves an item by its ID from a specified table.
- **UpdateItem:** Updates the Name attribute of an item in a table.
- **DeleteItem:** Deletes an item by its ID from a specified table.
- **DeleteTable:** Deletes a specified DynamoDB table.

### S3:
- **CreateBucket:** Creates an S3 bucket.
- **ListBuckets:** Lists all S3 buckets in the account.
- **UploadFile:** Uploads a file (represented as a string in this example) to the specified bucket with a specific object key.
- **DownloadFile:** Downloads a file from the specified bucket and prints its content. 
- **ListObjects:** Lists all objects within a specified bucket.
- **DeleteBucket:** Deletes a specified bucket. Calls deleteAllObjects to clear out all objects in the bucket before deletion, as S3 requires buckets to be empty before they can be deleted.
- **DeleteAllObjects:** Deletes all objects within a bucket, which is necessary before deleting the bucket itself.

### SNS
- **CreateTopic:** Creates an SNS topic and returns the topic ARN.
- **PublishMessage:** Publishes a message to the specified SNS topic.
- **ListTopics:** Lists all SNS topics in the account.
- **SubscribeToTopic:** Subscribes an endpoint (like email or HTTP) to a topic. This can be used to test various SNS subscription protocols (e.g., email, sms, http).
- **ListSubscriptions:** Lists all subscriptions for a specific topic by ARN. This can be useful for confirming active subscriptions.
- **DeleteTopic:** Deletes a specific SNS topic by its ARN, cleaning up after tests.

### SQS
- **SendMessage:** Sends a message to the specified queue URL. The function accepts the queueURL and message-body as parameters.
- **ReceiveMessages:** Receives messages from the specified queue URL using long polling (waits up to 5 seconds for a message).Prints out the message body and deletes the message after processing.
- **DeleteMessage:** Deletes a message from the queue using the queueURL and receiptHandle.This is called within ReceiveMessages to remove messages after they’ve been processed.
- **DeleteQueue:** Deletes the specified queue by queueURL. Used to clean up the queue after testing.

## Additional Configurations
- **Go Versions**: To test against multiple Go versions, use folders in the go/versions/ directory for each Go and AWS SDK version configuration.
- **Environment Variables**: 
    - Modify environment variables in the docker-compose.yml file for custom settings like AWS region or LocalStack endpoint.
    - Enable additional AWS services as needed by setting the `SERVICES` environment variable in LocalStack.

