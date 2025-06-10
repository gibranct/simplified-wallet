#!/bin/bash

QUEUE_URL="http://localhost:4566/000000000000/transaction-created-queue.fifo"
AWS_LOCAL="aws --endpoint-url=http://localhost:4566"

while true; do
  # Receive message
  MESSAGE=$($AWS_LOCAL sqs receive-message \
    --queue-url $QUEUE_URL \
    --attribute-names All \
    --message-attribute-names All \
    --max-number-of-messages 1 \
    --wait-time-seconds 20)

  # Check if a message was received
  if [ -n "$MESSAGE" ]; then
    # Extract receipt handle
    RECEIPT_HANDLE=$(echo $MESSAGE | jq -r '.Messages[0].ReceiptHandle')

    # Extract message body
    BODY=$(echo $MESSAGE | jq -r '.Messages[0].Body')

    # Process message (replace with your processing logic)
    echo "Processing message: $BODY"

    # Delete message
    $AWS_LOCAL sqs delete-message \
      --queue-url $QUEUE_URL \
      --receipt-handle "$RECEIPT_HANDLE"

    echo "Message processed and deleted"
  else
    echo "No messages available"
    sleep 5
  fi
done