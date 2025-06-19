package queue

import (
	"context"
	"fmt"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"
	"github.com/google/uuid"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SNSClient is an interface that defines the methods we use from the SNS client
// This allows us to mock the client for testing
type SNSClient interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

type SNS struct {
	client   SNSClient
	topicARN string
	otel     telemetry.Telemetry
}

// NewSNSWithClient creates a new SNS instance with a provided client
// This is primarily used for testing
func NewSNSWithClient(client SNSClient, otel telemetry.Telemetry) *SNS {
	// Default topic ARN for localstack (FIFO)
	topicARN := "arn:aws:sns:us-east-1:000000000000:transaction-events.fifo"
	if os.Getenv("SNS_TOPIC_ARN") != "" {
		topicARN = os.Getenv("SNS_TOPIC_ARN")
	}

	return &SNS{
		client:   client,
		topicARN: topicARN,
		otel:     otel,
	}
}

func NewSNS(otel telemetry.Telemetry) *SNS {
	// Default to localstack configuration
	endpoint := "http://localhost:4566"
	if os.Getenv("AWS_ENDPOINT") != "" {
		endpoint = os.Getenv("AWS_ENDPOINT")
	}

	// Default topic ARN for localstack (FIFO)
	topicARN := "arn:aws:sns:us-east-1:000000000000:transaction-events.fifo"
	if os.Getenv("SNS_TOPIC_ARN") != "" {
		topicARN = os.Getenv("SNS_TOPIC_ARN")
	}

	// Create custom AWS config for localstack
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           endpoint,
			SigningRegion: "us-east-1",
		}, nil
	})

	// Default AWS credentials for localstack
	awsAccessKeyID := "test"
	awsSecretAccessKey := "test"
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		awsAccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	}
	if os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}

	// Create AWS config
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	// Create SNS client with timeout
	client := sns.NewFromConfig(cfg, func(o *sns.Options) {
		o.RetryMaxAttempts = 3
		o.RetryMode = aws.RetryModeStandard
	})

	return &SNS{
		client:   client,
		topicARN: topicARN,
		otel:     otel,
	}
}

func (s *SNS) Send(ctx context.Context, message []byte) error {
	ctx, span := s.otel.Start(ctx, "SQS")
	defer span.End()

	// Log the message for debugging
	log.Printf("Sending message to SNS topic %s: %s", s.topicARN, string(message))

	// Create the publish input
	input := &sns.PublishInput{
		TopicArn: aws.String(s.topicARN),
		Message:  aws.String(string(message)),
		// FIFO-specific parameters
		MessageGroupId: aws.String("transaction-group"),
	}

	// Generate a unique deduplication ID based on the message content and current time
	// This is optional if ContentBasedDeduplication is enabled on the topic
	timestamp := time.Now().UnixNano()
	randID := uuid.New().String()
	input.MessageDeduplicationId = aws.String(fmt.Sprintf("%d-%s", timestamp, randID))

	// Publish the message
	_, err := s.client.Publish(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
