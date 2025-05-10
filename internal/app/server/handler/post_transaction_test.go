package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/app/server/handler"
	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostTransaction_InvalidJSONBody_ShouldReturn400(t *testing.T) {
	// Arrange
	createTransactionMock := &CreateTransactionMock{}
	handler := handler.New(createTransactionMock)

	// Create a mock HTTP request with invalid JSON
	reader := strings.NewReader(`{"amount": "not-a-number", "sender_id": "123"}`)
	r, _ := http.NewRequest("POST", "/transaction", reader)

	// Create a response recorder to record the HTTP response
	w := httptest.NewRecorder()

	// Act
	handler.PostTransaction(w, r)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var body map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Contains(t, body, "error")
	// The error message might vary, so we just check that there's an error message
	assert.NotEmpty(t, body["error"])
}

func TestPostTransaction_InvalidSenderID_ShouldReturn400(t *testing.T) {
	// Arrange
	createTransactionMock := &CreateTransactionMock{}
	handler := handler.New(createTransactionMock)

	// Create a mock HTTP request with invalid sender_id
	reader := strings.NewReader(`{"amount": 100, "sender_id": "invalid-uuid", "receiver_id": "89751234-abcd-1234-efgh-567890123456"}`)
	r, _ := http.NewRequest("POST", "/transaction", reader)

	// Create a response recorder to record the HTTP response
	w := httptest.NewRecorder()

	// Act
	handler.PostTransaction(w, r)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var body map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Contains(t, body, "error")
	assert.Equal(t, "invalid sender_id", body["error"])
}

func TestPostTransaction_InvalidReceiverID_ShouldReturn400(t *testing.T) {
	// Arrange
	createTransactionMock := &CreateTransactionMock{}
	handler := handler.New(createTransactionMock)

	// Create a mock HTTP request with invalid receiver_id
	// uuid
	reader := strings.NewReader(`{"amount": 100, "sender_id": "d6ae1675-5978-49d3-a6e3-619955ec6b2e", "receiver_id": "invalid-uuid"}`)
	r, _ := http.NewRequest("POST", "/transaction", reader)

	// Create a response recorder to record the HTTP response
	w := httptest.NewRecorder()

	// Act
	handler.PostTransaction(w, r)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var body map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Contains(t, body, "error")
	assert.Equal(t, "invalid receiver_id", body["error"])
}

func TestPostTransaction_WhenUsecaseReturnsError_ShouldReturn422(t *testing.T) {
	// Arrange
	expectedError := errors.New("transaction not allowed")

	createTransactionMock := &CreateTransactionMock{}
	createTransactionMock.On(
		"Execute",
		mock.Anything,
		mock.MatchedBy(func(input usecase.CreateTransactionInput) bool {
			return input.Amount == 100 &&
				input.SenderID.String() == "d6ae1675-5978-49d3-a6e3-619955ec6b2e" &&
				input.ReceiverID.String() == "f6de1685-5978-49d3-a6e3-619955ec6b2f"
		}),
	).Return("", expectedError)

	handler := handler.New(createTransactionMock)

	// Create a request with valid JSON data
	reqBody := `{
		"amount": 100, 
		"sender_id": "d6ae1675-5978-49d3-a6e3-619955ec6b2e", 
		"receiver_id": "f6de1685-5978-49d3-a6e3-619955ec6b2f"
	}`
	r, _ := http.NewRequest("POST", "/transaction", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	// Act
	handler.PostTransaction(w, r)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var body map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Contains(t, body, "error")
	assert.Equal(t, expectedError.Error(), body["error"])

	createTransactionMock.AssertExpectations(t)
}

func TestPostTransaction_WhenUsecaseSucceeds_ShouldReturn201WithTransactionID(t *testing.T) {
	// Arrange
	expectedTransactionID := "transaction-123"

	createTransactionMock := &CreateTransactionMock{}
	createTransactionMock.On(
		"Execute",
		mock.Anything,
		mock.Anything,
	).Return(expectedTransactionID, nil)

	handler := handler.New(createTransactionMock)

	// Create a request with valid JSON data
	reqBody := `{
		"amount": 100, 
		"sender_id": "d6ae1675-5978-49d3-a6e3-619955ec6b2e", 
		"receiver_id": "f6de1685-5978-49d3-a6e3-619955ec6b2f"
	}`
	r, _ := http.NewRequest("POST", "/transaction", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	// Act
	handler.PostTransaction(w, r)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var body map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Contains(t, body, "transaction_id")
	assert.Equal(t, expectedTransactionID, body["transaction_id"])

	createTransactionMock.AssertExpectations(t)
}

func TestPostTransaction_NegativeAmount_ShouldReturn422(t *testing.T) {
	// Arrange
	expectedError := errors.New("amount must be positive")

	createTransactionMock := &CreateTransactionMock{}
	createTransactionMock.On(
		"Execute",
		mock.Anything,
		mock.MatchedBy(func(input usecase.CreateTransactionInput) bool {
			return input.Amount == -50 &&
				input.SenderID.String() == "d6ae1675-5978-49d3-a6e3-619955ec6b2e" &&
				input.ReceiverID.String() == "f6de1685-5978-49d3-a6e3-619955ec6b2f"
		}),
	).Return("", expectedError)

	handler := handler.New(createTransactionMock)

	// Create a request with negative amount
	reqBody := `{
		"amount": -50, 
		"sender_id": "d6ae1675-5978-49d3-a6e3-619955ec6b2e", 
		"receiver_id": "f6de1685-5978-49d3-a6e3-619955ec6b2f"
	}`
	r, _ := http.NewRequest("POST", "/transaction", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	// Act
	handler.PostTransaction(w, r)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var body map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Contains(t, body, "error")
	assert.Equal(t, expectedError.Error(), body["error"])

	createTransactionMock.AssertExpectations(t)
}

type CreateTransactionMock struct {
	mock.Mock
}

func (m *CreateTransactionMock) Execute(ctx context.Context, input usecase.CreateTransactionInput) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}
