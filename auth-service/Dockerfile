FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy workspace file
COPY go.work .

# Copy proto module
COPY proto/ ./proto/

# Copy auth-service module
COPY auth-service/go.mod auth-service/go.sum ./auth-service/
WORKDIR /app/auth-service
RUN go mod download

# Copy auth-service source
COPY auth-service/ ./

RUN go build -o auth-service ./cmd/server


FROM alpine:latest
RUN apk --no-cache add ca-certificates 

WORKDIR /root/

COPY --from=builder /app/auth-service/auth-service .

EXPOSE 50051

CMD ["./auth-service"]
