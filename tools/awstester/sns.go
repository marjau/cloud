package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

const snsSvcName = "SNS"

type snsTester struct {
	awsTester
	svc *sns.Client
}

func newSNSTester(cfg aws.Config) *snsTester {
	return &snsTester{
		awsTester: awsTester{
			cfg:    cfg,
			logger: newPrefixedLogger(snsSvcName),
		},
		svc: sns.NewFromConfig(cfg),
	}
}

func (st *snsTester) RunTests() error {
	st.logger.log("Running Testing...")
	topicName := "TestTopic"
	message := "Hello, this is a test message!"
	protocol := "email"               // Change to "http" or another protocol as needed
	endpoint := "example@example.com" // Replace with a real endpoint for testing

	topicArn, err := st.createTopic(topicName)
	if err != nil {
		return err
	}
	defer st.clean(topicArn)

	if err := st.subscribeToTopic(topicArn, protocol, endpoint); err != nil {
		return err
	}
	if err := st.listSubscriptions(topicArn); err != nil {
		return err
	}
	if err := st.publishMessage(topicArn, message); err != nil {
		return err
	}
	if err := st.listTopics(); err != nil {
		return err
	}
	if err := st.deleteTopic(topicArn); err != nil {
		return err
	}

	st.logger.log("Testing Completed")
	return nil
}

func (st *snsTester) clean(topicArn string) {
	st.logger.log("Start Cleaning Testing...")
	st.logger.log("Testing Cleaning Completed")
}

func (st *snsTester) createTopic(topicName string) (string, error) {
	st.logger.logf("Creating Topic: %v", topicName)

	result, err := st.svc.CreateTopic(context.TODO(), &sns.CreateTopicInput{
		Name: aws.String(topicName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create topic %v: %w", topicName, err)
	}

	topicArn := aws.ToString(result.TopicArn)
	st.logger.logf("Topic %v created successfully: %s", topicName, topicArn)
	return topicArn, nil
}

func (st *snsTester) publishMessage(topicArn, message string) error {
	st.logger.logf("Publishing Message to Topic: %v", topicArn)

	st.logger.logf("Sending message: %v", message)
	_, err := st.svc.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(message),
	})
	if err != nil {
		return fmt.Errorf("failed to publish message to topic %v: %w", topicArn, err)
	}

	st.logger.log("Message published successfully")
	return nil
}

func (st *snsTester) listTopics() error {
	st.logger.log("Listing Topics")

	result, err := st.svc.ListTopics(context.TODO(), &sns.ListTopicsInput{})
	if err != nil {
		return fmt.Errorf("failed to list topics: %w", err)
	}

	st.logger.log("Topics:")
	for i, topic := range result.Topics {
		st.logger.logf("  %d. %s", i+1, aws.ToString(topic.TopicArn))
	}
	return nil
}

func (st *snsTester) subscribeToTopic(topicArn, protocol, endpoint string) error {
	st.logger.logf("Subscribing to Topic: %v", topicArn)

	_, err := st.svc.Subscribe(context.TODO(), &sns.SubscribeInput{
		TopicArn: aws.String(topicArn),
		Protocol: aws.String(protocol), // e.g., "email", "sms", "http", etc.
		Endpoint: aws.String(endpoint), // e.g., email address or phone number
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %v: %w", topicArn, err)
	}

	st.logger.log("Subscription request sent successfully")
	return nil
}

func (st *snsTester) listSubscriptions(topicArn string) error {
	st.logger.logf("Listing Subscriptions for Topic: %v", topicArn)

	result, err := st.svc.ListSubscriptionsByTopic(context.TODO(), &sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(topicArn),
	})
	if err != nil {
		return fmt.Errorf("failed to list subscriptions for topic %v: %w", topicArn, err)
	}

	st.logger.log("Subscriptions:")
	for i, subscription := range result.Subscriptions {
		st.logger.logf("  %d. %s", i+1, aws.ToString(subscription.SubscriptionArn))
	}
	return nil
}

func (st *snsTester) deleteTopic(topicArn string) error {
	st.logger.logf("Deleting Topic: %v", topicArn)

	_, err := st.svc.DeleteTopic(context.TODO(), &sns.DeleteTopicInput{
		TopicArn: aws.String(topicArn),
	})
	if err != nil {
		return fmt.Errorf("failed to delete topic %v: %w", topicArn, err)
	}

	st.logger.log("Topic deleted successfully")
	return nil
}
