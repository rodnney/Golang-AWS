package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSClient struct {
	client *sqs.Client
}

type Message struct {
	Body          string
	ReceiptHandle string
}

func NewSQSClient(cfg aws.Config) *SQSClient {
	return &SQSClient{
		client: sqs.NewFromConfig(cfg),
	}
}

func (s *SQSClient) ReceiveMessages(ctx context.Context, queueURL string, maxMessages int32) ([]Message, error) {
	output, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     20,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to receive messages from SQS: %w", err)
	}

	var messages []Message
	for _, msg := range output.Messages {
		messages = append(messages, Message{
			Body:          *msg.Body,
			ReceiptHandle: *msg.ReceiptHandle,
		})
	}

	return messages, nil
}

func (s *SQSClient) DeleteMessage(ctx context.Context, queueURL, receiptHandle string) error {
	_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		return fmt.Errorf("failed to delete message from SQS: %w", err)
	}

	return nil
}
