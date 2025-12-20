FROM golang:1.25-alpine AS dev

ARG MAIN_PKG
ARG DEBUG_PORT

ENV MAIN_PKG=${MAIN_PKG}
ENV DEBUG_PORT=${DEBUG_PORT}

WORKDIR /app

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY ../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../server .

ENTRYPOINT ["air", \
  "--build.cmd=go build -gcflags='all=-N -l' -o ./build/server ${MAIN_PKG}", \
  "--build.bin=./build/server", \
  "--build.full_bin=dlv exec --accept-multiclient --headless --continue --log --listen=:${DEBUG_PORT} --api-version=2 ./build/server --"]

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../server .
ENV CGO_ENABLED=0
RUN go build -o server ./cmd/api/main.go

FROM scratch AS runner

WORKDIR /app

COPY --from=builder /app/server .

ENTRYPOINT [ "./server" ]