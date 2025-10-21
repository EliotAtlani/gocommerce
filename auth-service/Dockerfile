# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Copy proto files first
COPY proto/auth ./proto/auth
COPY proto/user ./proto/user

# Copy auth-service
COPY auth-service ./auth-service

WORKDIR /build/auth-service

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/server

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /auth-service .

EXPOSE 50051

CMD ["./auth-service"]
