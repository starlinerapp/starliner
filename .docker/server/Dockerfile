FROM golang:1.25-alpine AS dev

ARG MAIN_PKG
ARG DEBUG_PORT
ARG BINARY_NAME

ENV MAIN_PKG=${MAIN_PKG}
ENV DEBUG_PORT=${DEBUG_PORT}
ENV BINARY_NAME=${BINARY_NAME}

WORKDIR /app

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
ENV PATH=$PATH:/root/go/bin

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

COPY --from=builder /app/${BINARY_NAME} ./app
COPY ../../server/.env .env

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
ENV PATH=$PATH:/root/go/bin

ENTRYPOINT [ "./app" ]