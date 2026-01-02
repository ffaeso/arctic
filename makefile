# Development environment variables (safe for local dev)
export ARCTIC_SERVER_ADDR ?= 9726
export ARCTIC_LOG_LEVEL ?= debug
export ARCTIC_LOG_FORMAT ?= text
export ARCTIC_DATASOURCE_USERNAME ?= arctic
export ARCTIC_DATASOURCE_PASSWORD ?= arctic_dev
export ARCTIC_DATASOURCE_HOST ?= localhost
export ARCTIC_DATASOURCE_PORT ?= 5432
export ARCTIC_DATASOURCE_DBNAME ?= arctic_dev
export ARCTIC_DATASOURCE_SSLMODE ?= disable

.PHONY: dev
dev: deps-up
	go run ./cmd/arctic start

.PHONY: build
build:
	go build -o arctic ./cmd/arctic

.PHONY: run
run: build
	./arctic start

.PHONY: deps-up
deps-up:
	@docker compose -f build/docker-compose.dev.yml up -d
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 2

.PHONY: deps-down
deps-down:
	docker compose -f build/docker-compose.dev.yml down

.PHONY: deps-logs
deps-logs:
	docker compose -f build/docker-compose.dev.yml logs -f

# standard test with coverage
.PHONY: test
test:
	go test ./... -v -cover

# for seeing html visualisation of whats covered and whats not
.PHONY: cover
cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: clean
clean:
	rm -f arctic
	docker compose -f build/docker-compose.dev.yml down -v

.PHONY: help
help:
	@echo "Arctic Development Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make dev         - Start dependencies and run Arctic (hot path for development)"
	@echo "  make build       - Build Arctic binary"
	@echo "  make run         - Build and run Arctic binary"
	@echo "  make test        - Run tests with coverage"
	@echo "  make cover       - Generate HTML coverage report"
	@echo "  make deps-up     - Start dependencies (PostgreSQL) in Docker"
	@echo "  make deps-down   - Stop dependencies"
	@echo "  make deps-logs   - Show dependency logs"
	@echo "  make clean       - Clean build artifacts and stop all containers"
	@echo ""
	@echo "Environment variables can be overridden:"
	@echo "  ARCTIC_LOG_LEVEL=warn make dev"
