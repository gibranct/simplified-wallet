package queue_test

import (
	"context"
	"errors"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"
	"os"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/provider/queue"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock SNS client that implements the queue.SNSClient interface
type mockSNSClient struct {
	mock.Mock
}

func (m *mockSNSClient) Publish(ctx context.Context, params *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sns.PublishOutput), args.Error(1)
}

// Test NewSNS constructor with default values
func TestNewSNS_ShouldUseDefaultValues(t *testing.T) {
	// Arrange
	// Clear environment variables to ensure defaults are used
	os.Clearenv()

	// Act
	snsClient := queue.NewSNS(telemetry.NewMockTelemetry())

	// Assert
	assert.NotNil(t, snsClient, "SNS client should not be nil")
}

func TestNewSNS_ShouldUseEnvironmentVariables(t *testing.T) {
	// Arrange
	os.Clearenv()
	err := os.Setenv("AWS_ENDPOINT", "http://custom-endpoint:4566")
	assert.NoError(t, err)
	err = os.Setenv("SNS_TOPIC_ARN", "custom-topic-arn")
	assert.NoError(t, err)
	err = os.Setenv("AWS_ACCESS_KEY_ID", "custom-key")
	assert.NoError(t, err)
	err = os.Setenv("AWS_SECRET_ACCESS_KEY", "custom-secret")
	assert.NoError(t, err)

	// Act
	snsClient := queue.NewSNS(telemetry.NewMockTelemetry())

	// Assert
	assert.NotNil(t, snsClient, "SNS client should not be nil")
}

func TestSNS_Send_ShouldReturnNilWhenPublishSucceeds(t *testing.T) {
	// Arrange
	ctx := context.Background()
	message := []byte("test message")

	// Create a mock SNS client that will return success
	mockOutput := &sns.PublishOutput{}
	mockClient := &mockSNSClient{}
	mockClient.On("Publish", mock.Anything, mock.Anything).Return(mockOutput, nil)

	// Create SNS instance with the mock client
	snsInstance := queue.NewSNSWithClient(mockClient, telemetry.NewMockTelemetry())

	// Act
	err := snsInstance.Send(ctx, message)

	// Assert
	assert.NoError(t, err, "Send should not return an error when publish succeeds")
	mockClient.AssertCalled(t, "Publish", mock.Anything, mock.Anything)
}

func TestSNS_Send_ShouldReturnErrorWhenPublishFails(t *testing.T) {
	// Arrange
	ctx := context.Background()
	message := []byte("test message")
	expectedErr := errors.New("publish error")

	// Create a mock SNS client that will return an error
	mockClient := &mockSNSClient{}
	mockClient.On("Publish", mock.Anything, mock.Anything).Return(nil, expectedErr)

	// Create SNS instance with the mock client
	snsInstance := queue.NewSNSWithClient(mockClient, telemetry.NewMockTelemetry())

	// Act
	err := snsInstance.Send(ctx, message)

	// Assert
	assert.Error(t, err, "Send should return an error when publish fails")
	assert.Contains(t, err.Error(), "failed to publish message")
	mockClient.AssertCalled(t, "Publish", mock.Anything, mock.Anything)
}
