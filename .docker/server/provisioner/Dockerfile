FROM golang:1.26-alpine AS dev

ENV PATH=/go/bin:$PATH:/root/.pulumi/bin

WORKDIR /app

RUN apk add --no-cache bash curl ansible openssh-client
RUN curl -fsSL https://get.pulumi.com | sh

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
RUN go build -o provisioner ./cmd/provisioner/main.go

FROM golang:1.26-alpine AS runner

WORKDIR /app

ENV PATH=/root/.pulumi/bin:$PATH

RUN apk add --no-cache bash curl ansible openssh-client
RUN curl -fsSL https://get.pulumi.com | sh

COPY --from=builder /app/provisioner ./provisioner
COPY ../../../server/.env .env

ENTRYPOINT [ "./provisioner" ]