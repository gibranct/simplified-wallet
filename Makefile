
.PHONY: build clean test run run/docker migrate migrate/down lint fmt help docker-up docker-down

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Application name
APP_NAME=simplified-wallet

# Docker related variables
DOCKER_COMPOSE=docker compose

# Script related variables
aws-local=

# Default target
.DEFAULT_GOAL := help

## build: Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(GOBIN)/$(APP_NAME) ./cmd/main.go

## clean: Clean build files
clean:
	@echo "Cleaning build files..."
	@rm -rf $(GOBIN)/*

## test: Run tests
test/unit:
	@echo "Running tests..."
	@go test -v ./internal/...

## test/coverage: Run tests with coverage
test/coverage:
	@echo "Running tests with coverage..."
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

## test/clean: Clean test cache and run tests
test/clean:
	@echo "Cleaning test cache and running tests..."
	@go clean -testcache
	@go test -v -count=1 ./...

## test/integration: Run integration tests
test/integration:
	@echo "Running integration tests..."
	@go test -v -tags=integration ./tests/...

## run: Run the application locally
run:
	@echo "Running $(APP_NAME)..."
	@go run ./cmd/main.go

## run/docker: Run the application in Docker
run/docker:
	@echo "Running $(APP_NAME) in Docker..."
	@$(DOCKER_COMPOSE) up --build app

## migrate: Run database migrations
migrate:
	@echo "Running database migrations..."
	@go run ./cmd/migrate/main.go up

## migrate/down: Rollback database migrations
migrate/down:
	@echo "Rolling back database migrations..."
	@go run ./cmd/migrate/main.go down

## lint: Run linters
lint:
	@echo "Running linters..."
	@golangci-lint run ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## docker-up: Start all Docker containers
docker-up:
	@echo "Starting Docker containers..."
	@$(DOCKER_COMPOSE) up -d

## docker-down: Stop all Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	@$(DOCKER_COMPOSE) down

## create-queue: Create SNS topic and SQS subscription
create-queue:
	@echo "Create topic and queue..."
	@./localstack-init.sh

## help: Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^##' Makefile | sed -e 's/## //g' | column -t -s ':'
