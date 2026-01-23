FROM golang:1.25-alpine AS dev

ENV PATH=/go/bin:$PATH:/root/.pulumi/bin

WORKDIR /app

RUN apk add --no-cache bash curl ansible openssh-client
RUN curl -fsSL https://get.pulumi.com | sh

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY ../../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../../server .

ENTRYPOINT ["air", \
  "--build.cmd=go build -gcflags='all=-N -l' -o ./build/provisioner ./cmd/provisioner/main.go", \
  "--build.bin=./build/provisioner", \
  "--build.full_bin=dlv exec --accept-multiclient --headless --continue --log --listen=:40002 --api-version=2 ./build/provisioner --"]

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ../../../server/go.mod ../../server/go.sum ./
RUN go mod download

COPY ../../../server .
ENV CGO_ENABLED=0
RUN go build -o provisioner ./cmd/provisioner/main.go

FROM golang:1.25-alpine AS runner

WORKDIR /app

ENV PATH=/root/.pulumi/bin:$PATH

RUN apk add --no-cache bash curl ansible openssh-client
RUN curl -fsSL https://get.pulumi.com | sh

COPY --from=builder /app/provisioner ./provisioner

ENTRYPOINT [ "./provisioner" ]