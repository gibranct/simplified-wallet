
# Simplified Wallet

A simplified wallet application that allows users to create accounts, make transactions, and manage merchants. The application uses a microservices architecture with message queues for asynchronous processing.

## Technologies

- **Golang**: Main programming language
- **PostgreSQL**: Database for storing user and transaction data
- **AWS SNS/SQS**: Message queue for asynchronous transaction processing (using LocalStack for local development)
- **Docker & Docker Compose**: Containerization and orchestration

## Project Structure

```
simplified-wallet/
├── api/                  # API documentation and HTTP request examples
├── cmd/                  # Application entry points
│   ├── main.go           # Main application
│   └── migrate/          # Database migration tool
├── internal/             # Internal packages
│   ├── app/              # Application logic
│   │   ├── server/       # HTTP server and handlers
│   │   └── usecase/      # Business logic use cases
│   ├── domain/           # Domain models and business rules
│   └── provider/         # External service providers (DB, queue, etc.)
├── migrations/           # Database migration files
├── tests/                # Test files
├── docker-compose.yml    # Docker Compose configuration
├── Makefile              # Build and run commands
├── localstack-init.sh    # LocalStack initialization script
└── consume_sqs_messages.sh # Script to consume SQS messages
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.16 or higher
- AWS CLI (for local testing with LocalStack)

### Running the Application

1. Start the required services (PostgreSQL and LocalStack):

```shell
make docker-up
```

2. Run the application:

```shell
make run
```

3. Create SNS topic and SQS queue:

```shell
make create-queue
```

Alternatively, you can run the application in Docker:

```shell
make run/docker
```

### Database Migrations

Run database migrations:

```shell
make migrate
```

Rollback migrations:

```shell
make migrate/down
```

## API Endpoints

### Create Common User

```http
POST /v1/users HTTP/1.1
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@mail.com",
  "password": "securepassword123",
  "cpf": "12345678901"
}
```

### Create Merchant User

```http
POST /v1/merchants HTTP/1.1
Content-Type: application/json

{
  "name": "Business Corp",
  "email": "business@corp.com",
  "password": "securepassword123",
  "cnpj": "12345678000190"
}
```

### Create Transaction

```http
POST /v1/transactions HTTP/1.1
Content-Type: application/json

{
  "amount": 100.99,
  "sender_id": "7250961f-c104-46dd-9447-d57b4f5a2be4",
  "receiver_id": "d47d6618-7f43-47dc-a33c-be833f5e6ef8"
}
```

## Message Processing

The application uses AWS SNS and SQS (via LocalStack for local development) for asynchronous transaction processing.

### Consuming SQS Messages

You can use the provided script to consume messages from the SQS queue:

```shell
./consume_sqs_messages.sh
```

## Development

### Running Tests

Run unit tests:

```shell
make test/unit
```

Run integration tests:

```shell
make test/integration
```

Run tests with coverage:

```shell
make test/coverage
```

### Code Quality

Format code:

```shell
make fmt
```

Run linters:

```shell
make lint
```

## Help

For a list of all available commands:

```shell
make help
```

## Metrics and Tracing

The application includes metrics and distributed tracing capabilities:

### Prometheus Metrics

Metrics are exposed at the `/metrics` endpoint and can be visualized in the Prometheus dashboard:

```shell
make metrics
```

Key metrics include:
- `http_requests_total`: Total number of HTTP requests
- `http_request_duration_seconds`: Duration of HTTP requests
- `transactions_total`: Total number of transactions
- `transactions_amount_total`: Total amount of transactions

### Jaeger Tracing

Distributed tracing is implemented using Jaeger. View traces in the Jaeger UI:

```shell
make tracing
```

This provides insights into:
- Request flow through the system
- Performance bottlenecks
- Error propagation