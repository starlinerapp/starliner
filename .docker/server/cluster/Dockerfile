FROM golang:1.25-alpine AS dev

ENV PATH=/go/bin:$PATH

WORKDIR /app

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY ../../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../../server .

ENTRYPOINT ["air", \
  "--build.cmd=go build -gcflags='all=-N -l' -o ./build/cluster ./cmd/cluster/main.go", \
  "--build.bin=./build/cluster", \
  "--build.delay=1000", \
  "--build.send_interrupt=true", \
  "--build.kill_delay=2000", \
  "--build.full_bin=dlv exec --accept-multiclient --headless --continue --log --listen=:40002 --api-version=2 ./build/cluster --" \
]

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ../../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../../server .
ENV CGO_ENABLED=0
RUN go build -o cluster ./cmd/cluster/main.go

FROM golang:1.25-alpine AS runner

WORKDIR /app

COPY --from=builder /app/cluster ./cluster
COPY ../../../server/.env .env

ENTRYPOINT [ "./cluster" ]