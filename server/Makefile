SHELL := /bin/bash

.PHONY: deps deps-cleancache swag proto format lint fix security upgrade help

GOCMD=go

deps: ## Install dependencies
	$(GOCMD) get -u -t -v ./...
	$(GOCMD) mod tidy

deps-cleancache: ## Clear cache in Go module
	$(GOCMD) clean -modcache

swag: ## Generate swagger docs
	swag init -g ./cmd/api/main.go -o ./cmd/api/docs

proto: ## Generate Go code from protobufs
	protoc \
		--go_out=. --go_opt=paths=source_relative $(shell find ./internal/core/infrastructure/grpc/proto/v1 -name "*.proto") \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative $(shell find ./internal/core/infrastructure/grpc/proto/v1 -name "*.proto")

format:
	gofmt -w .
	find . -name "*.sql" -exec pg_format -c .pg_format -i {} \;

lint:
	@golangci-lint run

fix:
	@golangci-lint run --fix

security:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

upgrade:
	go get -u ./...
	go mod tidy -v
	go list -m -u all

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'