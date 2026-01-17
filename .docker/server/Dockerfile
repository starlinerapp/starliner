FROM golang:1.25-alpine AS builder

ARG MAIN_PKG
ARG BINARY_NAME

WORKDIR /app

COPY ../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../server .
ENV CGO_ENABLED=0
RUN go build -o ${BINARY_NAME} ${MAIN_PKG}

FROM golang:1.25-alpine AS runner

ARG BINARY_NAME

WORKDIR /app

ENV PATH=/root/.pulumi/bin:$PATH

RUN apk add --no-cache bash curl ansible openssh-client
RUN curl -fsSL https://get.pulumi.com | sh

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY --from=builder /app/${BINARY_NAME} ./app
COPY ../../server/.env .env
COPY ../../server/internal/infrastructure/postgres/migrations ./internal/infrastructure/postgres/migrations

ENTRYPOINT [ "./app" ]