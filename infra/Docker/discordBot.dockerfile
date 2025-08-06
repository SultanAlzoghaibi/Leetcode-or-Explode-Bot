# ----- Stage 1: Build -----
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go/go.mod ./
COPY go/go.sum ./
RUN go mod download

COPY go/. ./

# Build the Discord bot binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o discordBot ./cmd/discordBot

# ----- Stage 2: Run -----
FROM alpine:latest
WORKDIR /app/

COPY --from=builder /app/discordBot .
RUN chmod +x discordBot

COPY go/credentials.json ./credentials.json

EXPOSE 9100
CMD ["./discordBot"]