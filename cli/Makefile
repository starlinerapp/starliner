SHELL := /bin/bash

.PHONY: deps deps-cleancache lint fix security upgrade build openapi-fetch openapi-generate

GOCMD=go

OPENAPI_URL=http://localhost:5173/api/auth/open-api/generate-schema
OPENAPI_SPEC=internal/infrastructure/auth/generated/openapi.json

deps: ## Install dependencies
	$(GOCMD) get -u -t -v ./...
	$(GOCMD) mod tidy

deps-cleancache: ## Clear cache in Go module
	$(GOCMD) clean -modcache

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

build:
	$(GOCMD) build -o build/starliner ./cmd/main.go

openapi-generate:
	mkdir -p internal/infrastructure/auth/generated/client
	curl -s $(OPENAPI_URL) > $(OPENAPI_SPEC)
	openapi-generator-cli generate --skip-validate-spec -i $(OPENAPI_SPEC) -g go -o internal/infrastructure/auth/generated/client
