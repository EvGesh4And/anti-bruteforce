# ───── Go ─────
generate:          ## go generate ./...
	go generate ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.8

lint: install-lint-deps
	go mod download
	golangci-lint run ./...

test:
	go test -race ./internal/...

.PHONY: generate lint test
