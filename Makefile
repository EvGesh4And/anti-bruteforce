# ───── Настройки сборки ─────
BIN          ?= ./bin
TAG          ?= develop
SERVICE_NAME ?= antibruteforce

GIT_HASH     := $(shell git rev-parse --short HEAD)
BUILD_DATE   := $(shell date -u +%Y-%m-%dT%H:%M:%S)
LDFLAGS      := -X main.release=$(TAG) \
                -X main.buildDate=$(BUILD_DATE) \
                -X main.gitHash=$(GIT_HASH)

DOCKERFILE   := build/Dockerfile
DOCKER_COMPOSE := deployments/docker-compose.yaml

# ───── Утилиты ─────
generate:          ## go generate ./...
	go generate ./...

install-lint-deps: ## установить golangci-lint при необходимости
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "→ Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.8; \
	fi

lint: install-lint-deps ## линтинг
	go mod download
	golangci-lint run ./...

test: ## юнит-тесты
	go test -race ./internal/...

# ───── CLI ─────
$(BIN):
	mkdir -p $(BIN)

cli: $(BIN) ## сборка CLI
	go build -o $(BIN)/cli ./cmd/cli
	$(BIN)/cli -h

# ───── Docker ─────
build: ## сборка Docker-образа сервера
	docker build \
		-t $(SERVICE_NAME):$(TAG) \
		--build-arg LDFLAGS='$(LDFLAGS)' \
		--build-arg SERVICE_NAME=$(SERVICE_NAME) \
		-f $(DOCKERFILE) .

up: build ## поднятие Docker Compose
	docker compose -f $(DOCKER_COMPOSE) up -d

down: ## остановить Docker Compose
	docker compose -f $(DOCKER_COMPOSE) down

# ───── Общее ─────
.PHONY: generate install-lint-deps lint test cli build up down
