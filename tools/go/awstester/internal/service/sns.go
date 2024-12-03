package service

import (
	"context"
	"fmt"
	"myapps/awstester/internal/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

const snsSvcName = "SNS"

type SnsTester struct {
	awsTester
	svc *sns.Client
}

func NewSNSTester(cfg aws.Config) *SnsTester {
	return &SnsTester{
		awsTester: awsTester{
			cfg:    cfg,
			logger: logger.NewPrefixedLogger(snsSvcName),
		},
		svc: sns.NewFromConfig(cfg),
	}
}

func (st SnsTester) GetName() string {
	return snsSvcName
}

func (st SnsTester) Run() error {
	st.logger.Log("Running Testing...")
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

	st.logger.Log("Testing Completed")
	return nil
}

func (st SnsTester) clean(topicArn string) {
	st.logger.Log("Start Cleaning Testing...")
	st.logger.Log("Testing Cleaning Completed")
}

func (st SnsTester) createTopic(topicName string) (string, error) {
	st.logger.Logf("Creating Topic: %v", topicName)

	result, err := st.svc.CreateTopic(context.TODO(), &sns.CreateTopicInput{
		Name: aws.String(topicName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create topic %v: %w", topicName, err)
	}

	topicArn := aws.ToString(result.TopicArn)
	st.logger.Logf("Topic %v created successfully: %s", topicName, topicArn)
	return topicArn, nil
}

func (st SnsTester) publishMessage(topicArn, message string) error {
	st.logger.Logf("Publishing Message to Topic: %v", topicArn)

	st.logger.Logf("Sending message: %v", message)
	_, err := st.svc.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(message),
	})
	if err != nil {
		return fmt.Errorf("failed to publish message to topic %v: %w", topicArn, err)
	}

	st.logger.Log("Message published successfully")
	return nil
}

func (st SnsTester) listTopics() error {
	st.logger.Log("Listing Topics")

	result, err := st.svc.ListTopics(context.TODO(), &sns.ListTopicsInput{})
	if err != nil {
		return fmt.Errorf("failed to list topics: %w", err)
	}

	st.logger.Log("Topics:")
	for i, topic := range result.Topics {
		st.logger.Logf("  %d. %s", i+1, aws.ToString(topic.TopicArn))
	}
	return nil
}

func (st SnsTester) subscribeToTopic(topicArn, protocol, endpoint string) error {
	st.logger.Logf("Subscribing to Topic: %v", topicArn)

	_, err := st.svc.Subscribe(context.TODO(), &sns.SubscribeInput{
		TopicArn: aws.String(topicArn),
		Protocol: aws.String(protocol), // e.g., "email", "sms", "http", etc.
		Endpoint: aws.String(endpoint), // e.g., email address or phone number
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %v: %w", topicArn, err)
	}

	st.logger.Log("Subscription request sent successfully")
	return nil
}

func (st SnsTester) listSubscriptions(topicArn string) error {
	st.logger.Logf("Listing Subscriptions for Topic: %v", topicArn)

	result, err := st.svc.ListSubscriptionsByTopic(context.TODO(), &sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(topicArn),
	})
	if err != nil {
		return fmt.Errorf("failed to list subscriptions for topic %v: %w", topicArn, err)
	}

	st.logger.Log("Subscriptions:")
	for i, subscription := range result.Subscriptions {
		st.logger.Logf("  %d. %s", i+1, aws.ToString(subscription.SubscriptionArn))
	}
	return nil
}

func (st SnsTester) deleteTopic(topicArn string) error {
	st.logger.Logf("Deleting Topic: %v", topicArn)

	_, err := st.svc.DeleteTopic(context.TODO(), &sns.DeleteTopicInput{
		TopicArn: aws.String(topicArn),
	})
	if err != nil {
		return fmt.Errorf("failed to delete topic %v: %w", topicArn, err)
	}

	st.logger.Log("Topic deleted successfully")
	return nil
}
