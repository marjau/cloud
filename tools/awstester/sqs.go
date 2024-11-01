package main

// TODO: implement testing.T approach

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const sqsSvcName = "SQS"

type sqsTester struct {
	awsTester
	svc *sqs.Client
}

func newSQSTester(cfg aws.Config) *sqsTester {
	sqs.NewFromConfig(cfg)
	return &sqsTester{
		awsTester: awsTester{
			cfg:    cfg,
			logger: newPrefixedLogger(sqsSvcName),
		},
		svc: sqs.NewFromConfig(cfg),
	}
}

func (st sqsTester) RunTests() error {
	st.logger.log("Start Testing...")
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

	st.logger.log("Testing Completed.")
	return nil
}

func (st sqsTester) clean(queueURL string) {
	st.logger.log("Start Cleaning Testing...")
	if err := st.deleteQueue(queueURL); err != nil {
		if strings.Contains(err.Error(), "AWS.SimpleQueueService.NonExistentQueue") {
			st.logger.logf("INFO: Queue not found: %v", queueURL)
		} else {
			st.logger.logf("ERROR: %v", err)
		}
	}
	st.logger.log("Cleaning Testing Completed")
}

func (st sqsTester) listQueues() error {
	st.logger.log("List Queues")

	result, err := st.svc.ListQueues(context.TODO(), &sqs.ListQueuesInput{})
	if err != nil {
		return fmt.Errorf("failed to list queues: %w", err)
	}

	st.logger.log("Queues:")
	for i, queueUrl := range result.QueueUrls {
		st.logger.logf("  %d. %s", i+1, aws.ToString(&queueUrl))
	}
	return nil
}

func (st *sqsTester) createQueue(queueName string) (string, error) {
	st.logger.logf("Creating Queue: %v", queueName)

	result, err := st.svc.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create %v queue: %w", queueName, err)
	}

	st.logger.logf("Queue %v created successfully: %s", queueName, aws.ToString(result.QueueUrl))
	return aws.ToString(result.QueueUrl), nil
}

func (st *sqsTester) sendMessage(queueURL, messageBody string) error {
	st.logger.logf("Sending Message to Queue: %v", queueURL)

	st.logger.logf("Sending message: %v", messageBody)
	_, err := st.svc.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(messageBody),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to queue %v: %w", queueURL, err)
	}

	st.logger.log("Message sent successfully")
	return nil
}

func (st *sqsTester) receiveMessages(queueURL string) error {
	st.logger.logf("Receiving Messages from Queue: %v", queueURL)

	result, err := st.svc.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: 1, // Number of messages to receive
		WaitTimeSeconds:     5, // Long polling for up to 5 seconds
	})
	if err != nil {
		return fmt.Errorf("failed to receive messages from queue %v: %w", queueURL, err)
	}

	if len(result.Messages) == 0 {
		st.logger.log("No messages received")
		return nil
	}

	// Process and delete received messages
	for _, message := range result.Messages {
		st.logger.logf("Received message: %s", aws.ToString(message.Body))
		if err := st.deleteMessage(queueURL, aws.ToString(message.ReceiptHandle)); err != nil {
			return err
		}
	}
	return nil
}

func (st *sqsTester) deleteMessage(queueURL, receiptHandle string) error {
	st.logger.logf("Deleting Message from Queue: %v", queueURL)

	_, err := st.svc.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		return fmt.Errorf("failed to delete message from queue %v: %w", queueURL, err)
	}

	st.logger.log("Message deleted successfully")
	return nil
}

func (st *sqsTester) deleteQueue(queueURL string) error {
	st.logger.logf("Deleting Queue: %v", queueURL)

	_, err := st.svc.DeleteQueue(context.TODO(), &sqs.DeleteQueueInput{
		QueueUrl: aws.String(queueURL),
	})
	if err != nil {
		return fmt.Errorf("failed to delete queue %v: %w", queueURL, err)
	}

	st.logger.log("Queue deleted successfully")
	return nil
}
