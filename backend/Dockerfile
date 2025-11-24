# Base stage
FROM golang:1.23-alpine AS base
WORKDIR /app
COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

# Development stage
FROM base AS dev
RUN go install github.com/cosmtrek/air@v1.49.0
COPY . .
CMD ["air", "-c", ".air.toml"]

# Builder stage
FROM base AS builder
COPY . .
RUN go build -o main ./cmd/server/main.go

# Production stage
FROM alpine:latest AS prod
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
