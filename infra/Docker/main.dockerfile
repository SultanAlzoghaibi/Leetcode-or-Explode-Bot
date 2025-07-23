# ----- Stage 1: Build -----
FROM golang:1.24-alpine AS builder
RUN pwd
WORKDIR /go
RUN pwd
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# ----- Stage 2: Run -----
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/main .

CMD ["./main"]