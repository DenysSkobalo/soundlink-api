SERVER_NAME := soundlink-api
MAIN := ./cmd/server/main.go
BUILD_DIR := ./bin
GO := go
ENV_FILE := .env

GO_BUILD_FLAGS := 

.PHONY: all
all: run

.PHONY: build
build:
	@echo "Building $(SERVER_NAME) in $(ENV) mode..."
	$(GO) build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(SERVER_NAME) $(MAIN)

.PHONY: run
run:
	$(GO) run $(MAIN)

.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test -v ./...

.PHONY: clean
clean:
	@echo "Cleaning build files..."
	rm -rf $(BUILD_DIR)

