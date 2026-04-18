# Binary name
BINARY_NAME=go-ipc
# Output directory
BIN_DIR=bin
# Main package path
MAIN_PATH=./cmd

all: build

## build: Build the binary
build:
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PATH)

## run: Build and run the application
run: build
	@./$(BIN_DIR)/$(BINARY_NAME)

## clean: Remove built binaries
clean:
	go clean
	rm -rf $(BIN_DIR)

## race: Detect race condition
race:
	go run -race ./...

## test: Run tests including race condition identifications
test:
	@echo "Running tests"
	@go test -v ./...
	@echo "Running tests to identify race condition"
	@go test -race ./...


## fmt: Format code
fmt:
	go fmt ./...


## vet: Run go vet
vet:
	go vet ./...

## vuln: Check for known vulnerabilities
vuln:
	@echo "Running vulnerability check..."
	@govulncheck ./... || echo "⚠️ Vulnerabilities found"

## install-tools: Install development tools
install-tools:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go get github.com/spf13/viper
