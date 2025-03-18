APP_NAME=chatrbox
BUILD_DIR=build
MAIN=./cmd/main.go

.PHONY: all test build fmt lint clean run

all: build

# Format code
fmt:
	go fmt ./...

# Run tests
test:
	go test -v ./...

# Build binary
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)

# Run the app
run:
	go run $(MAIN)

# Lint (requires golangci-lint installed)
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)