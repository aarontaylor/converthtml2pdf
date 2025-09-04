.PHONY: build build-all test install clean release lint fmt

BINARY_NAME=converthtml2pdf
VERSION=1.0.0
GOFLAGS=-ldflags="-s -w -X 'main.version=$(VERSION)'"

build:
	go build $(GOFLAGS) -o $(BINARY_NAME) main.go

build-all: build-linux build-darwin build-windows

build-linux:
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(BINARY_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -o $(BINARY_NAME)-linux-arm64 main.go

build-darwin:
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o $(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o $(BINARY_NAME)-darwin-arm64 main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o $(BINARY_NAME)-windows-amd64.exe main.go

test:
	go test -v -cover ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

install: build
	sudo mv $(BINARY_NAME) /usr/local/bin/

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html

release: clean build-all
	mkdir -p releases
	mv $(BINARY_NAME)-* releases/

run:
	go run main.go

fmt:
	go fmt ./...

lint:
	go vet ./...
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

deps:
	go mod download
	go mod tidy

help:
	@echo "Available targets:"
	@echo "  build        - Build for current platform"
	@echo "  build-all    - Build for all platforms"
	@echo "  test         - Run tests"
	@echo "  install      - Install locally"
	@echo "  clean        - Clean build artifacts"
	@echo "  release      - Create release binaries"
	@echo "  run          - Run the application"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linters"
	@echo "  deps         - Download and tidy dependencies"