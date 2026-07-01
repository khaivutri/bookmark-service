.PHONY: help deps tidy fmt vet test test-coverage build run clean all

GO ?= go
MAIN_PATH ?= ./cmd/api
BIN_DIR ?= ./bin
BINARY ?= $(BIN_DIR)/bookmark-service
APP_PORT ?= 8080
SERVICE_NAME ?= bookmark_service
INSTANCE_ID ?=

help:
	@echo "Available targets:"
	@echo "  make help              Show this help message"
	@echo "  make deps              Download Go modules"
	@echo "  make tidy              Clean up module dependencies"
	@echo "  make fmt               Format Go source files"
	@echo "  make vet               Run go vet"
	@echo "  make test              Run all unit and integration tests"
	@echo "  make test-coverage     Run tests with coverage report"
	@echo "  make build             Build the binary into $(BIN_DIR)"
	@echo "  make run               Run the application"
	@echo "  make clean             Remove build artifacts and coverage files"
	@echo "  make all               Run fmt, vet, test, and build"

deps:
	$(GO) mod download

tidy:
	$(GO) mod tidy

fmt:
	gofmt -w ./cmd ./internal

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

test-coverage:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BINARY) $(MAIN_PATH)

run:
	APP_PORT=$(APP_PORT) SERVICE_NAME=$(SERVICE_NAME) INSTANCE_ID=$(INSTANCE_ID) $(GO) run $(MAIN_PATH)

clean:
	rm -rf $(BIN_DIR) coverage.out coverage.html

all: fmt vet test build
