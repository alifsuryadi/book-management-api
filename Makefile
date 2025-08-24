# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=book-management-api

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
run:
	$(GOCMD) run main.go

# Build for Linux (for deployment)
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v ./...

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code
vet:
	$(GOCMD) vet ./...

# Install development dependencies
install-deps:
	$(GOGET) github.com/gin-gonic/gin@latest
	$(GOGET) github.com/lib/pq@latest
	$(GOGET) github.com/rubenv/sql-migrate@latest
	$(GOGET) github.com/golang-jwt/jwt/v5@latest
	$(GOGET) golang.org/x/crypto@latest

# Docker commands
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)

# Railway commands
railway-login:
	railway login

railway-deploy:
	railway up

# Development workflow
dev: deps fmt vet test run

# Production build
prod: clean build-linux

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  run          - Run the application"
	@echo "  build-linux  - Build for Linux deployment"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  install-deps - Install development dependencies"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  railway-login - Login to Railway"
	@echo "  railway-deploy - Deploy to Railway"
	@echo "  dev          - Development workflow (deps, fmt, vet, test, run)"
	@echo "  prod         - Production build"
	@echo "  help         - Show this help message"

.PHONY: build clean test deps run build-linux fmt vet install-deps docker-build docker-run railway-login railway-deploy dev prod help