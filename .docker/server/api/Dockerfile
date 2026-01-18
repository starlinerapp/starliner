FROM golang:1.25-alpine AS dev

WORKDIR /app

ENV PATH=/go/bin:$PATH

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY ../../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../../server .

ENTRYPOINT ["air", \
  "--build.cmd=go build -gcflags='all=-N -l' -o ./build/api ./cmd/api/main.go", \
  "--build.bin=./build/api", \
  "--build.full_bin=dlv exec --accept-multiclient --headless --continue --log --listen=:40000 --api-version=2 ./build/api --"]

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ../../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../../server .
ENV CGO_ENABLED=0
RUN go build -o api ./cmd/api/main.go

FROM golang:1.25-alpine AS runner

WORKDIR /app

COPY --from=builder /app/api ./api
COPY ../../../server/.env .env
COPY ../../../server/internal/core/infrastructure/postgres/migrations ./internal/core/infrastructure/postgres/migrations
COPY --from=dev /go/bin/goose /usr/local/bin/goose

ENTRYPOINT [ "./api" ]