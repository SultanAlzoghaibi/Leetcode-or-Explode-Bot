# ----- Stage 1: Build -----
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go/go.mod ./
COPY go/go.sum ./
RUN go mod download

COPY go/. ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

# ----- Stage 2: Run -----
FROM alpine:latest
WORKDIR /app/

COPY --from=builder /app/main .
COPY go/credentials.json ./credentials.json

EXPOSE 9100

CMD ["./main"]