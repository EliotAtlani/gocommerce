# Multi-stage build for smaller final image
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go workspace files first (for dependency caching)
COPY go.work go.work.sum ./

# Copy all necessary modules (workspace requires all modules in go.work)
COPY api-gateway/ ./api-gateway/
COPY auth-service/ ./auth-service/
COPY user-service/ ./user-service/
COPY proto/ ./proto/

# Build the application
WORKDIR /build/api-gateway
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway ./cmd/server

# Final stage - minimal runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /build/api-gateway/gateway .

# Expose HTTP port
EXPOSE 8080

# Run the gateway
CMD ["./gateway"]
