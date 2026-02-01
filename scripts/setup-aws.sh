#!/bin/bash

# AWS Setup Script for LocalStack
# Configures SNS, SQS, DynamoDB and S3 resources

export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
ENDPOINT_URL="http://localhost:4566"

echo "Configuring LocalStack resources..."

# 1. Create SNS Topic
echo "Creating SNS topic..."
aws --endpoint-url=$ENDPOINT_URL sns create-topic --name transaction-events

# 2. Create SQS Queues
echo "Creating SQS queues..."
aws --endpoint-url=$ENDPOINT_URL sqs create-queue --queue-name validation-queue
aws --endpoint-url=$ENDPOINT_URL sqs create-queue --queue-name enrichment-queue
aws --endpoint-url=$ENDPOINT_URL sqs create-queue --queue-name audit-queue
aws --endpoint-url=$ENDPOINT_URL sqs create-queue --queue-name audit-queue-dlq

# 3. Subscribe Queues to SNS Topic
TOPIC_ARN="arn:aws:sns:us-east-1:000000000000:transaction-events"

aws --endpoint-url=$ENDPOINT_URL sns subscribe \
    --topic-arn $TOPIC_ARN \
    --protocol sqs \
    --notification-endpoint http://localhost:4566/000000000000/validation-queue

aws --endpoint-url=$ENDPOINT_URL sns subscribe \
    --topic-arn $TOPIC_ARN \
    --protocol sqs \
    --notification-endpoint http://localhost:4566/000000000000/enrichment-queue

aws --endpoint-url=$ENDPOINT_URL sns subscribe \
    --topic-arn $TOPIC_ARN \
    --protocol sqs \
    --notification-endpoint http://localhost:4566/000000000000/audit-queue

# 4. Create DynamoDB Table
echo "Creating DynamoDB table..."
aws --endpoint-url=$ENDPOINT_URL dynamodb create-table \
    --table-name transactions \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5

# 5. Create S3 Bucket
echo "Creating S3 bucket..."
aws --endpoint-url=$ENDPOINT_URL s3 mb s3://audit-logs

echo "AWS resources configured successfully!"
