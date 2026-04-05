# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go
# Redis only (for local `make run` against Docker Redis)
docker-run:
	@if docker compose version >/dev/null 2>&1; then \
		docker compose up -d redis_bp; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose up -d redis_bp; \
	else \
		echo "Docker Compose not found. Install the Docker Compose V2 plugin (docker compose)."; \
		exit 1; \
	fi

# API + Redis (builds image from Dockerfile)
docker-up:
	@if docker compose version >/dev/null 2>&1; then \
		docker compose up --build; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose up --build; \
	else \
		echo "Docker Compose not found. Install the Docker Compose V2 plugin (docker compose)."; \
		exit 1; \
	fi

# Shutdown Compose services
docker-down:
	@if docker compose version >/dev/null 2>&1; then \
		docker compose down; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose down; \
	else \
		echo "Docker Compose not found. Install the Docker Compose V2 plugin (docker compose)."; \
		exit 1; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch docker-run docker-up docker-down itest
