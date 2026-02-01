package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSClient struct {
	client *sns.Client
}

func NewSNSClient(cfg aws.Config) *SNSClient {
	return &SNSClient{
		client: sns.NewFromConfig(cfg),
	}
}

func (s *SNSClient) Publish(ctx context.Context, topicARN string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = s.client.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(string(body)),
		TopicArn: aws.String(topicARN),
	})
	if err != nil {
		return fmt.Errorf("failed to publish to SNS: %w", err)
	}

	return nil
}
