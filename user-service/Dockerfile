# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Copy proto files first
COPY proto/user ./proto/user

# Copy user-service
COPY user-service ./user-service

WORKDIR /build/user-service

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ./cmd/server

# Run stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /user-service .

# Copy migrations
COPY --from=builder /build/user-service/migrations ./migrations

EXPOSE 50052

CMD ["./user-service"]
