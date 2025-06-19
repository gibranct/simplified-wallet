package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostMerchant_Integration(t *testing.T) {
	// Create a request with valid JSON data
	reqBody := `{
        "name": "John Doe",
        "email": "john@example.com",
        "password": "password123",
        "cnpj": "46797901000157"
    }`
	resp, err := http.Post(server.URL+"/v1/merchants", "application/json", bytes.NewBufferString(reqBody))
	require.NoError(t, err)
	defer func() {
		err = resp.Body.Close()
		require.NoError(t, err)
	}()

	// Check the status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Parse the response body
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	// Check that a user_id was returned
	userID, ok := response["user_id"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, userID)
}
