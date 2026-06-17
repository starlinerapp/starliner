FROM golang:1.26-alpine AS dev

WORKDIR /app

ENV PATH=/go/bin:$PATH

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY ../../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../../server .

ENTRYPOINT ["air", "-c", "/config/air.toml"]

FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY ../../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../../server .
ENV CGO_ENABLED=0
RUN go build -o builder ./cmd/builder/main.go

FROM golang:1.26-alpine AS runner

WORKDIR /app

COPY --from=builder /app/builder ./builder
COPY ../../../server/.env .env

ENTRYPOINT [ "./builder" ]