FROM golang:1.25-alpine AS dev

WORKDIR /app

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENTRYPOINT ["air", \
  "--build.cmd=go build -gcflags='all=-N -l' -o ./build/server ./cmd/api/main.go", \
  "--build.bin=./build/server", \
  "--build.full_bin=dlv exec --accept-multiclient --headless --continue --log --listen=:40000 --api-version=2 ./build/server --"]

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -o server ./cmd/api/main.go

FROM scratch AS runner

WORKDIR /app

COPY --from=builder /app/server .

ENTRYPOINT [ "./server" ]