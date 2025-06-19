package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/app/server/handler"
	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostUser_InvalidJSONBody_ShouldReturn400(t *testing.T) {
	// Arrange
	createUserMock := &CreateUserMock{}
	createTransactionMock := &CreateTransactionMock{}
	mockTelemetry := telemetry.NewMockTelemetry()
	h := handler.New(createTransactionMock, createUserMock, mockTelemetry)

	// Create a mock HTTP request with invalid JSON
	reader := strings.NewReader(`{"name": "John Doe", "email": "invalid-json`)
	r, _ := http.NewRequest("POST", "/v1/users", reader)

	// Create a response recorder to record the HTTP response
	w := httptest.NewRecorder()

	// Act
	h.PostUser(w, r)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var body map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Contains(t, body, "error")
	assert.NotEmpty(t, body["error"])

	createUserMock.AssertNotCalled(t, "Execute")
}

func TestPostUser_EmailAlreadyRegistered_ShouldReturn422(t *testing.T) {
	// Arrange
	expectedError := errors.New("email already registered")

	createUserMock := &CreateUserMock{}
	createTransactionMock := &CreateTransactionMock{}
	createUserMock.On(
		"Execute",
		mock.Anything,
		mock.MatchedBy(func(input usecase.CreateUserInput) bool {
			return input.Name == "John Doe" &&
				input.Email == "john@example.com" &&
				input.Password == "password123" &&
				input.Document == "12345678901"
		}),
	).Return("", expectedError)

	mockTelemetry := telemetry.NewMockTelemetry()
	h := handler.New(createTransactionMock, createUserMock, mockTelemetry)

	// Create a request with valid JSON data
	reqBody := `{
		"name": "John Doe",
		"email": "john@example.com",
		"password": "password123",
		"cpf": "12345678901"
	}`
	r, _ := http.NewRequest("POST", "/v1/users", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	// Act
	h.PostUser(w, r)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var body map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Contains(t, body, "error")
	assert.Equal(t, expectedError.Error(), body["error"])

	createUserMock.AssertExpectations(t)
}

func TestPostUser_SuccessfulCreation_ShouldReturn201WithUserID(t *testing.T) {
	// Arrange
	expectedUserID := "user-123"

	createUserMock := &CreateUserMock{}
	createTransactionMock := &CreateTransactionMock{}
	createUserMock.On(
		"Execute",
		mock.Anything,
		mock.MatchedBy(func(input usecase.CreateUserInput) bool {
			return input.Name == "John Doe" &&
				input.Email == "john@example.com" &&
				input.Password == "password123" &&
				input.Document == "12345678901"
		}),
	).Return(expectedUserID, nil)

	mockTelemetry := telemetry.NewMockTelemetry()
	h := handler.New(createTransactionMock, createUserMock, mockTelemetry)

	// Create a request with valid JSON data
	reqBody := `{
		"name": "John Doe",
		"email": "john@example.com",
		"password": "password123",
		"cpf": "12345678901"
	}`
	r, _ := http.NewRequest("POST", "/v1/users", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	// Act
	h.PostUser(w, r)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var body map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Contains(t, body, "user_id")
	assert.Equal(t, expectedUserID, body["user_id"])

	createUserMock.AssertExpectations(t)
}
