# Development environment variables (safe for local dev)
#
# NOTE: These environment variables using makefile has 
# precendence over /etc/arctic/config.yml or ~/.arctic/config.yml.
# If running development environment using makefile, 
# make sure to modify configurations down below instead of 
# using config.yml
export ARCTIC_SERVER_ADDR ?= 9726
export ARCTIC_LOG_LEVEL ?= debug
export ARCTIC_LOG_FORMAT ?= text
export ARCTIC_DATASOURCE_USERNAME ?= arctic
export ARCTIC_DATASOURCE_PASSWORD ?= arctic
export ARCTIC_DATASOURCE_HOST ?= localhost
export ARCTIC_DATASOURCE_PORT ?= 5432
export ARCTIC_DATASOURCE_DBNAME ?= arctic_db
export ARCTIC_DATASOURCE_SSLMODE ?= disable

# hot path -> spins up dependencies (postgres) via 
# docker compose + runs arctic start
.PHONY: dev
dev: deps-up
	go run ./cmd/arctic start

# builds binary into ./bin/arctic
.PHONY: build
build:
	go build -o ./bin/arctic ./cmd/arctic

# runs the binary built in make build
.PHONY: run
run: build
	./bin/arctic start

# spins up dependencies
.PHONY: deps-up
deps-up:
	@docker compose -f build/docker-compose.dev.yml up -d
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 2

# stops the running dependencies
.PHONY: deps-down
deps-down:
	docker compose -f build/docker-compose.dev.yml down

# run this to see logs coming from docker compose 
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

# stops running containers, removes it, and remove the volume
.PHONY: clean
clean:
	rm -f bin/
	docker compose -f build/docker-compose.dev.yml down -v

# help menu
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
