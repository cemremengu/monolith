SHELL := bash

.DEFAULT_GOAL := help

PACKAGE ?= monolith

VERSION ?= $(shell git describe --tags 2>/dev/null || echo "v0.0.0")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE_BUILT ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

CACHE_DIR ?= $(CURDIR)/.cache
GOCACHE ?= $(CACHE_DIR)/go-build

GO_FLAGS ?=
TAGS ?=
LD_FLAGS ?=
CGO_ENABLED ?= 0

GO_LD_FLAGS := $(LD_FLAGS) -X $(PACKAGE).Version=$(VERSION) -X $(PACKAGE).Commit=$(COMMIT) -X $(PACKAGE).DateBuilt=$(DATE_BUILT)

.PHONY: help
help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make <target>\n\nTargets:\n"} /^[a-zA-Z0-9_.-]+:.*##/ {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: go-cache
go-cache:
	@mkdir -p "$(GOCACHE)"

.PHONY: run
run: ## Run backend + frontend dev servers (parallel)
	@$(MAKE) -j2 server web

.PHONY: server
server: ## Run backend with hot reload (wgo)
	wgo run ./cmd/monolith

.PHONY: web
web: ## Run frontend dev server
	cd web && npm run dev

.PHONY: lint
lint: ## Run Go + web linting
	@$(MAKE) lint-go
	@$(MAKE) lint-web

.PHONY: lint-go
lint-go: ## Run Go formatting + linting (golangci-lint)
	golangci-lint fmt
	golangci-lint run

.PHONY: lint-web
lint-web: ## Run web linting + typecheck
	cd web && npm run lint:fix
	cd web && npm run check-types

.PHONY: router-generate
router-generate: ## Generate TanStack Router route tree
	cd web && npm run router:generate

.PHONY: update
update: ## Update Go dependencies and tidy modules
	GOCACHE=$(GOCACHE) go get -u ./...
	GOCACHE=$(GOCACHE) go mod tidy

.PHONY: test
test: test-go ## Run Go tests

.PHONY: test-go
test-go: go-cache ## Run Go tests
	GOCACHE=$(GOCACHE) go test -v ./...

.PHONY: test-web
test-web: ## Run web tests
	cd web && npm test

.PHONY: test-coverage
test-coverage: go-cache ## Run Go tests with coverage
	GOCACHE=$(GOCACHE) go test -v -coverprofile=coverage.out ./...
	GOCACHE=$(GOCACHE) go tool cover -html=coverage.out -o coverage.html

.PHONY: test-all
test-all: ## Run all tests
	@$(MAKE) test-go
	@$(MAKE) test-web

.PHONY: build-web
build-web: ## Build the web frontend
	cd web && npm run build

.PHONY: build-server
build-server: ## Build the Go server binary (host OS/arch)
	@$(MAKE) go-cache
	GOCACHE=$(GOCACHE) CGO_ENABLED=$(CGO_ENABLED) go build $(GO_FLAGS) -tags "$(TAGS)" -ldflags "$(GO_LD_FLAGS)" ./cmd/monolith

.PHONY: build-win
build-win: ## Build for Windows (amd64)
	@$(MAKE) build-web
	@$(MAKE) go-cache
	GOCACHE=$(GOCACHE) GOOS=windows GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(GO_FLAGS) -tags "$(TAGS)" -ldflags "$(GO_LD_FLAGS)" ./cmd/monolith

.PHONY: build-linux
build-linux: ## Build for Linux (amd64)
	@$(MAKE) build-web
	@$(MAKE) go-cache
	GOCACHE=$(GOCACHE) GOOS=linux GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(GO_FLAGS) -tags "$(TAGS)" -ldflags "$(GO_LD_FLAGS)" ./cmd/monolith

.PHONY: build-docker
build-docker: ## Build Docker image
	docker build -t monolith:latest .

.PHONY: build-all
build-all: ## Build for Windows, Linux, and Docker
	@$(MAKE) build-win
	@$(MAKE) build-linux
	@$(MAKE) build-docker

.PHONY: install-tools
install-tools: ## Install development tools (golangci-lint, wgo)
	GOCACHE=$(GOCACHE) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	GOCACHE=$(GOCACHE) go install github.com/wgo/wgo@latest

.PHONY: install-web
install-web: ## Install web dependencies
	cd web && npm install

.PHONY: setup
setup: ## Setup development environment
	@$(MAKE) install-tools
	@$(MAKE) install-web

.PHONY: clean
clean: ## Clean build artifacts and caches
	rm -rf dist/ web/dist/
	rm -f coverage.out coverage.html
	rm -rf "$(CACHE_DIR)"
	GOCACHE=$(GOCACHE) go clean -cache
