#!/bin/bash

# Wait for LocalStack to be ready
echo "Waiting for LocalStack to be ready..."
while ! nc -z localstack 4566; do
  sleep 1
done
echo "LocalStack is ready!"

# Create SNS topic (FIFO)
echo "Creating FIFO SNS topic..."
awslocal sns create-topic --name transaction-events.fifo --attributes FifoTopic=true,ContentBasedDeduplication=true

# Create SQS queue (FIFO)
echo "Creating FIFO SQS queue..."
awslocal sqs create-queue --queue-name transaction-created-queue.fifo --attributes FifoQueue=true,ContentBasedDeduplication=true

# Subscribe SQS queue to SNS topic
echo "Subscribing SQS queue to SNS topic..."
TOPIC_ARN=$(awslocal sns list-topics --query 'Topics[0].TopicArn' --output text)
QUEUE_URL=$(awslocal sqs get-queue-url --queue-name transaction-created-queue.fifo --query 'QueueUrl' --output text)
QUEUE_ARN=$(awslocal sqs get-queue-attributes --queue-url $QUEUE_URL --attribute-names QueueArn --query 'Attributes.QueueArn' --output text)

awslocal sns subscribe \
  --topic-arn $TOPIC_ARN \
  --protocol sqs \
  --notification-endpoint $QUEUE_ARN

echo "Initialization complete!"
