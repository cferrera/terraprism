.PHONY: build test clean install lint run help docker-build docker-build-mac

# Build variables
BINARY_NAME=terraprism
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR=bin
# uname -m spellings differ from Go's; normalise for the docker cross-build targets
HOST_ARCH?=$(shell uname -m | sed 's/x86_64/amd64/; s/aarch64/arm64/')
GO_FILES=$(shell find . -name '*.go' -type f)
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

# Default target
all: build

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: Build the binary
build: $(BUILD_DIR)/$(BINARY_NAME)

$(BUILD_DIR)/$(BINARY_NAME): $(GO_FILES)
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/terraprism

## install: Install the binary to $GOPATH/bin
install:
	go install $(LDFLAGS) ./cmd/terraprism

## test: Run tests
test:
	go test -v -race -cover ./...

## test-coverage: Run tests with coverage report
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## lint: Run linter
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## fmt: Format code
fmt:
	go fmt ./...
	gofmt -s -w .

## clean: Remove build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## run: Run the application (reads from stdin)
run: build
	@echo "Run with: terraform plan -no-color | ./$(BUILD_DIR)/$(BINARY_NAME)"

## demo: Run with sample plan
demo: build
	@echo "Terraform will perform the following actions:" | cat - testdata/sample-plan.txt 2>/dev/null | ./$(BUILD_DIR)/$(BINARY_NAME) || \
		echo "No sample plan found. Create testdata/sample-plan.txt or pipe real terraform output."

## docker-build: Build a Linux binary in Docker (no local Go toolchain needed)
docker-build:
	docker build -o dist --build-arg VERSION=$(VERSION) .
	@file dist/terraprism

## docker-build-mac: Build a macOS binary in Docker for this host's arch
docker-build-mac:
	docker build -o dist --build-arg VERSION=$(VERSION) \
		--build-arg GOOS=darwin --build-arg GOARCH=$(HOST_ARCH) .
	@file dist/terraprism

## deps: Download dependencies
deps:
	go mod download
	go mod tidy

## release: Build for multiple platforms
release: clean
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/terraprism
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/terraprism
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/terraprism
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/terraprism
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/terraprism
	@echo "Binaries built in $(BUILD_DIR)/"
	@ls -la $(BUILD_DIR)/

