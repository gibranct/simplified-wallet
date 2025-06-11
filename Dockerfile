# Build stage
FROM golang:1.23.3-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/simplified-wallet ./cmd/main.go

# Final stage
FROM alpine:3.18

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/simplified-wallet .
# Copy migrations folder for database migrations
COPY --from=builder /app/migrations ./migrations

# Expose the application port
EXPOSE 3000

# Set environment variables
ENV POSTGRES_URI=postgres://postgres:postgres@postgres:5432/wallet
ENV POSTGRES_SSLMODE=false

# Run the application
CMD ["./simplified-wallet"]