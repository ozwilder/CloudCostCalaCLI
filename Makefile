.PHONY: help build run test clean fmt lint

help:
	@echo "CloudCostCalaCLI - Cloud Asset Inventory Tool"
	@echo ""
	@echo "Available targets:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application with default config"
	@echo "  make test        - Run tests"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Run linter"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make all         - Build and run"

build:
	go build -o bin/cloudcostcala ./cmd/cloudcostcala

run: build
	./bin/cloudcostcala --config config.example.json --output cloud-assets-inventory.xlsx

test:
	go test -v ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/
	rm -f cloud-assets-inventory.xlsx

all: clean build run
