package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

const (
	testTopicName = "test"
)

type SNSTester struct {
	tester
	svc *sns.SNS
}

func NewSNSTester(cfg AWSConfig, sess *session.Session) *SNSTester {
	return &SNSTester{
		tester: tester{
			config: cfg,
			logger: newPrefixedLogger("SNS"),
		},
		svc: sns.New(sess),
	}
}

func (st SNSTester) RunTests() error {
	st.logger.log("Start testing...")
	if err := st.createTopic(); err != nil {
		return err
	}
	if err := st.listTopics(); err != nil {
		return err
	}
	st.logger.log("Testing completed")
	return nil
}

func (st SNSTester) Clean() {
	st.logger.log("Cleaning tests...")
	st.logger.log("Tests cleaned")
}

func (st SNSTester) listTopics() error {
	st.logger.log("List Topics")

	result, err := st.svc.ListTopics(&sns.ListTopicsInput{})
	if err != nil {
		return err
	}

	st.logger.log("Topics:")
	for i, topic := range result.Topics {
		st.logger.logf("  %d. %s", i+1, aws.StringValue(topic.TopicArn))
	}

	return nil
}

func (st SNSTester) createTopic() error {
	st.logger.logf("Create %v topic", testTopicName)

	_, err := st.svc.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String(testTopicName),
	})
	if err != nil {
		return err
	}
	st.logger.logf("Topic %v created successfully", testTopicName)
	return nil
}
