package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSTester struct {
	tester
	svc *sqs.SQS
}

func NewSQSTester(cfg AWSConfig, sess *session.Session) *SQSTester {
	return &SQSTester{
		tester: tester{
			config: cfg,
			logger: newPrefixedLogger("SQS"),
		},
		svc: sqs.New(sess),
	}
}

func (st SQSTester) RunTests() error {
	st.logger.log("Start testing...")
	if err := st.CreateQueue(); err != nil {
		return err
	}
	if err := st.ListQueues(); err != nil {
		return err
	}
	st.logger.log("Testing completed")
	return nil
}

func (st SQSTester) Clean() {
	st.logger.log("Cleaning tests...")
	st.logger.log("Tests cleaned")
}

func (st SQSTester) ListQueues() error {
	st.logger.log("List Queues")

	result, err := st.svc.ListQueues(&sqs.ListQueuesInput{})
	if err != nil {
		return err
	}

	st.logger.log("Queues:")
	for i, queueUrl := range result.QueueUrls {
		st.logger.logf("  %d %s", i+1, aws.StringValue(queueUrl))
	}
	return nil
}

func (st SQSTester) CreateQueue() error {
	qName := "TestQueue"
	st.logger.logf("Create %v Queue", qName)

	createQueueInput := &sqs.CreateQueueInput{
		QueueName: aws.String(qName),
	}
	result, err := st.svc.CreateQueue(createQueueInput)
	if err != nil {
		return err
	}

	st.logger.logf("Queue %v created successfully: %s", qName, aws.StringValue(result.QueueUrl))
	return nil
}
