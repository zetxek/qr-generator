# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Only copy go.mod and go.sum first, to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the source code
COPY . .

# Build the binary
RUN go build -o server ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
