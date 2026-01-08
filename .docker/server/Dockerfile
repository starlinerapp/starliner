FROM golang:1.25-alpine AS dev

ARG MAIN_PKG
ARG DEBUG_PORT
ARG BINARY_NAME

ENV MAIN_PKG=${MAIN_PKG}
ENV DEBUG_PORT=${DEBUG_PORT}
ENV BINARY_NAME=${BINARY_NAME}
ENV PATH=/go/bin:$PATH:/root/.pulumi/bin

WORKDIR /app

RUN apk add --no-cache bash curl ansible openssh-client
RUN curl -fsSL https://get.pulumi.com | sh

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY ../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../server .

ENTRYPOINT ["air", \
  "--build.cmd=go build -gcflags='all=-N -l' -o ./build/${BINARY_NAME} ${MAIN_PKG}", \
  "--build.bin=./build/${BINARY_NAME}", \
  "--build.full_bin=dlv exec --accept-multiclient --headless --continue --log --listen=:${DEBUG_PORT} --api-version=2 ./build/${BINARY_NAME} --"]

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

COPY --from=builder /app/${BINARY_NAME} ./app
COPY ../../server/.env .env
COPY ../../server/internal/infrastructure/postgres/migrations ./internal/infrastructure/postgres/migrations
COPY --from=dev /go/bin/goose /usr/local/bin/goose

ENTRYPOINT [ "./app" ]