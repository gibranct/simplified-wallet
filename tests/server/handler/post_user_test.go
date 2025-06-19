package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com.br/gibranct/simplified-wallet/internal/app/server/router"
	test "github.com.br/gibranct/simplified-wallet/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var server *httptest.Server

func TestMain(m *testing.M) {
	ctx := context.Background()
	migrateVersion := uint(2) // Use the latest migration version

	// Setup
	container, db, err := test.SetupTestDatabase(ctx, migrateVersion)
	if err != nil {
		panic(err)
	}
	otel, err := telemetry.NewJaeger(context.Background(), "")
	if err != nil {
		panic(err)
	}
	defer func() {
		err = container.Terminate(ctx)
		if err != nil {
			panic(err)
		}
	}()
	defer func() {
		err = db.Close()
		if err != nil {
			panic(err)
		}
	}()
	defer func() {
		err = otel.Shutdown(ctx)
		if err != nil {
			panic(err)
		}
	}()

	r := router.InitRoutes(otel)

	server = httptest.NewServer(r)
	defer server.Close()
	code := m.Run() // Run the tests
	os.Exit(code)
}

func TestPostUser_Integration(t *testing.T) {
	// Create a request with valid JSON data
	reqBody := `{
        "name": "John Doe 2",
        "email": "john2@example.com",
        "password": "password123",
        "cpf": "12345678901"
    }`
	resp, err := http.Post(server.URL+"/v1/users", "application/json", bytes.NewBufferString(reqBody))
	require.NoError(t, err)
	defer func() {
		err = resp.Body.Close()
		require.NoError(t, err)
	}()

	// Check the status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Parse the response body
	var jsonResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	require.NoError(t, err)

	// Check that a user_id was returned
	userID, ok := jsonResponse["user_id"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, userID)
}
