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