package service

// TODO: implement testing.T approach

import (
	"context"
	"fmt"
	"strings"

	"myapps/awstester/internal/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const sqsSvcName = "SQS"

type SqsTester struct {
	awsTester
	svc *sqs.Client
}

func NewSQSTester(cfg aws.Config) *SqsTester {
	sqs.NewFromConfig(cfg)
	return &SqsTester{
		awsTester: awsTester{
			cfg:    cfg,
			logger: logger.NewPrefixedLogger(sqsSvcName),
		},
		svc: sqs.NewFromConfig(cfg),
	}
}

func (st SqsTester) GetName() string {
	return sqsSvcName
}

func (st SqsTester) Run() error {
	st.logger.Log("Start Testing...")
	queueName := "TestQueue"
	message := "Hello, this is a test message!"

	queueURL, err := st.createQueue(queueName)
	if err != nil {
		return err
	}
	defer st.clean(queueURL)

	if err := st.listQueues(); err != nil {
		return err
	}
	if err := st.sendMessage(queueURL, message); err != nil {
		return err
	}
	if err := st.receiveMessages(queueURL); err != nil {
		return err
	}
	if err := st.deleteQueue(queueURL); err != nil {
		return err
	}

	st.logger.Log("Testing Completed.")
	return nil
}

func (st SqsTester) clean(queueURL string) {
	st.logger.Log("Start Cleaning Testing...")
	if err := st.deleteQueue(queueURL); err != nil {
		if strings.Contains(err.Error(), "AWS.SimpleQueueService.NonExistentQueue") {
			st.logger.Logf("INFO: Queue not found: %v", queueURL)
		} else {
			st.logger.Logf("ERROR: %v", err)
		}
	}
	st.logger.Log("Cleaning Testing Completed")
}

func (st SqsTester) listQueues() error {
	st.logger.Log("List Queues")

	result, err := st.svc.ListQueues(context.TODO(), &sqs.ListQueuesInput{})
	if err != nil {
		return fmt.Errorf("failed to list queues: %w", err)
	}

	st.logger.Log("Queues:")
	for i, queueUrl := range result.QueueUrls {
		st.logger.Logf("  %d. %s", i+1, aws.ToString(&queueUrl))
	}
	return nil
}

func (st SqsTester) createQueue(queueName string) (string, error) {
	st.logger.Logf("Creating Queue: %v", queueName)

	result, err := st.svc.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create %v queue: %w", queueName, err)
	}

	st.logger.Logf("Queue %v created successfully: %s", queueName, aws.ToString(result.QueueUrl))
	return aws.ToString(result.QueueUrl), nil
}

func (st SqsTester) sendMessage(queueURL, messageBody string) error {
	st.logger.Logf("Sending Message to Queue: %v", queueURL)

	st.logger.Logf("Sending message: %v", messageBody)
	_, err := st.svc.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(messageBody),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to queue %v: %w", queueURL, err)
	}

	st.logger.Log("Message sent successfully")
	return nil
}

func (st SqsTester) receiveMessages(queueURL string) error {
	st.logger.Logf("Receiving Messages from Queue: %v", queueURL)

	result, err := st.svc.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: 1, // Number of messages to receive
		WaitTimeSeconds:     5, // Long polling for up to 5 seconds
	})
	if err != nil {
		return fmt.Errorf("failed to receive messages from queue %v: %w", queueURL, err)
	}

	if len(result.Messages) == 0 {
		st.logger.Log("No messages received")
		return nil
	}

	// Process and delete received messages
	for _, message := range result.Messages {
		st.logger.Logf("Received message: %s", aws.ToString(message.Body))
		if err := st.deleteMessage(queueURL, aws.ToString(message.ReceiptHandle)); err != nil {
			return err
		}
	}
	return nil
}

func (st SqsTester) deleteMessage(queueURL, receiptHandle string) error {
	st.logger.Logf("Deleting Message from Queue: %v", queueURL)

	_, err := st.svc.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		return fmt.Errorf("failed to delete message from queue %v: %w", queueURL, err)
	}

	st.logger.Log("Message deleted successfully")
	return nil
}

func (st SqsTester) deleteQueue(queueURL string) error {
	st.logger.Logf("Deleting Queue: %v", queueURL)

	_, err := st.svc.DeleteQueue(context.TODO(), &sqs.DeleteQueueInput{
		QueueUrl: aws.String(queueURL),
	})
	if err != nil {
		return fmt.Errorf("failed to delete queue %v: %w", queueURL, err)
	}

	st.logger.Log("Queue deleted successfully")
	return nil
}
