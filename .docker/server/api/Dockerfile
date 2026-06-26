FROM golang:1.26-alpine AS dev

WORKDIR /app

ENV PATH=/go/bin:$PATH

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY ../../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../../server .

RUN go build -o /usr/local/bin/starliner ./cmd/cli

ENTRYPOINT ["air", "-c", "/config/air.toml"]

FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY ../../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../../server .
ENV CGO_ENABLED=0
RUN go build -o api ./cmd/api/main.go && \
    go build -o starliner ./cmd/cli

FROM golang:1.26-alpine AS runner

WORKDIR /app

COPY --from=builder /app/api ./api
COPY --from=builder /app/starliner /usr/local/bin/starliner
COPY ../../../server/.env .env
COPY ../../../server/internal/api/infrastructure/postgres/migrations ./internal/api/infrastructure/postgres/migrations
COPY --from=dev /go/bin/goose /usr/local/bin/goose

ENTRYPOINT [ "./api" ]