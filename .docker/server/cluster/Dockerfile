FROM golang:1.26-alpine AS dev

ENV PATH=/go/bin:$PATH

WORKDIR /app

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
RUN go build -o cluster ./cmd/cluster/main.go

FROM golang:1.26-alpine AS runner

WORKDIR /app

COPY --from=builder /app/cluster ./cluster
COPY ../../../server/.env .env

ENTRYPOINT [ "./cluster" ]