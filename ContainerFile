# Stage 1: Build the Go application
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /getMe .

# Stage 2: Create the final lightweight image
FROM alpine:latest

COPY --from=builder /getMe /getMe

VOLUME /data/getMeStore

ENTRYPOINT ["/getMe"]
