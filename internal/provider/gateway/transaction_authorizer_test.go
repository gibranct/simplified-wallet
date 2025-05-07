package gateway_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/provider/gateway"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionAuthorizer_IsTransactionAllowed_ShouldReturnFalseWhenHTTPRequestCreationFails(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	// Cancel the context to force NewRequestWithContext to fail
	cancel()

	mockHttpClient := &http.Client{}
	authorizer := gateway.NewTransactionAuthorizer(mockHttpClient)

	// Act
	result := authorizer.IsTransactionAllowed(ctx)

	// Assert
	assert.False(t, result, "Should return false when HTTP request creation fails")
}

func TestTransactionAuthorizer_IsTransactionAllowed_ShouldReturnFalseWhenHTTPClientDoFails(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Create a mock HTTP client that will return an error on Do
	mockHttpClient := &mockHTTPClient{}
	mockHttpClient.On("Do", mock.Anything).Return(nil, errors.New("connection error"))

	authorizer := gateway.NewTransactionAuthorizer(mockHttpClient)

	// Act
	result := authorizer.IsTransactionAllowed(ctx)

	// Assert
	assert.False(t, result, "Should return false when HTTP client Do fails")
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

func TestTransactionAuthorizer_IsTransactionAllowed_ShouldReturnTrueWhenResponseStatusCodeIsOK(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Create a mock HTTP client that will return a 200 OK response
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       http.NoBody,
	}
	mockHttpClient := &mockHTTPClient{}
	mockHttpClient.On("Do", mock.Anything).Return(mockResp, nil)

	authorizer := gateway.NewTransactionAuthorizer(mockHttpClient)

	// Act
	result := authorizer.IsTransactionAllowed(ctx)

	// Assert
	assert.True(t, result, "Should return true when API response status code is 200 OK")
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

func TestTransactionAuthorizer_IsTransactionAllowed_ShouldReturnFalseWhenResponseStatusCodeIsNotOK(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Create a mock HTTP client that will return a non-200 response
	mockResp := &http.Response{
		StatusCode: http.StatusForbidden,
		Body:       http.NoBody,
	}
	mockHttpClient := &mockHTTPClient{}
	mockHttpClient.On("Do", mock.Anything).Return(mockResp, nil)

	authorizer := gateway.NewTransactionAuthorizer(mockHttpClient)

	// Act
	result := authorizer.IsTransactionAllowed(ctx)

	// Assert
	assert.False(t, result, "Should return false when API response status code is not 200 OK")
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

// Mock HTTP client
type mockHTTPClient struct {
	mock.Mock
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}
