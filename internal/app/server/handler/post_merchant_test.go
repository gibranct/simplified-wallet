package handler_test

import (
	"encoding/json"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/app/server/handler"
	"github.com.br/gibranct/simplified-wallet/internal/app/usecase"
	"github.com.br/gibranct/simplified-wallet/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostMerchant_EmptyRequestBody_ShouldReturn400(t *testing.T) {
	// Arrange
	createUserMock := &CreateUserMock{}
	h := handler.New(nil, createUserMock, telemetry.NewMockTelemetry())

	// Create a mock HTTP request with an empty body
	r, _ := http.NewRequest("POST", "/v1/merchants", strings.NewReader(""))

	// Create a response recorder to record the HTTP response
	w := httptest.NewRecorder()

	// Act
	h.PostMerchant(w, r)

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

func TestPostMerchant_CreateUserUseCaseSuccess_ShouldReturn201Created(t *testing.T) {
	// Arrange
	createUserMock := &CreateUserMock{}
	h := handler.New(nil, createUserMock, telemetry.NewMockTelemetry())

	// Create a mock HTTP request with valid input
	reader := strings.NewReader(`{"name": "John Doe", "email": "johndoe@example.com", "password": "password123", "cnpj": "47775767000156"}`)
	r, _ := http.NewRequest("POST", "/v1/merchants", reader)

	// Create a response recorder to record the HTTP response
	w := httptest.NewRecorder()

	// Mock the use case to return a successful user ID
	expectedUserID := "d6ae1675-5978-49d3-a6e3-619955ec6b2e"
	createUserMock.On(
		"Execute",
		mock.Anything,
		mock.MatchedBy(func(input usecase.CreateUserInput) bool {
			return input.Name == "John Doe" &&
				input.Email == "johndoe@example.com" &&
				input.Password == "password123" &&
				input.Document == "47775767000156" &&
				input.UserType == vo.MerchantUserType
		}),
	).Return(expectedUserID, nil)

	// Act
	h.PostMerchant(w, r)

	// Assert
	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var body map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Contains(t, body, "user_id")
	assert.Equal(t, expectedUserID, body["user_id"])
}
