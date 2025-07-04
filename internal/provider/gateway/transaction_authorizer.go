package gateway

import (
	"context"
	"github.com.br/gibranct/simplified-wallet/internal/provider/telemetry"
	"log"
	"net/http"
	"os"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type TransactionAuthorizer struct {
	httpClient Client
	logger     *log.Logger
	otel       telemetry.Telemetry
}

func (ta *TransactionAuthorizer) IsTransactionAllowed(ctx context.Context) bool {
	const url = "https://util.devi.tools/api/v2/authorize"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		ta.logger.Printf("Error creating HTTP request: %v\n", err)
		return false
	}
	resp, err := ta.httpClient.Do(req)
	if err != nil {
		return false
	}
	return resp.StatusCode == http.StatusOK
}

func NewTransactionAuthorizer(httpClient Client, otel telemetry.Telemetry) *TransactionAuthorizer {
	return &TransactionAuthorizer{
		httpClient: httpClient,
		logger:     log.New(os.Stdout, "transaction_authorizer: ", log.LstdFlags),
		otel:       otel,
	}
}
