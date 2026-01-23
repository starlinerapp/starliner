FROM golang:1.25-alpine AS dev

WORKDIR /app

ENV PATH=/go/bin:$PATH

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY ../../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../../server .

ENTRYPOINT ["air", \
  "--build.cmd=go build -gcflags='all=-N -l' -o ./build/builder ./cmd/builder/main.go", \
  "--build.bin=./build/builder", \
  "--build.full_bin=dlv exec --accept-multiclient --headless --continue --log --listen=:40001 --api-version=2 ./build/builder --"]

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ../../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../../server .
ENV CGO_ENABLED=0
RUN go build -o builder ./cmd/builder/main.go

FROM golang:1.25-alpine AS runner

WORKDIR /app

COPY --from=builder /app/builder ./builder
COPY ../../../server/.env .env

ENTRYPOINT [ "./builder" ]