# API Gateway - GoCommerce

## Overview

The API Gateway is the **single entry point** for all client requests in the GoCommerce microservices architecture. It acts as a reverse proxy that translates REST/JSON requests from clients into gRPC calls to backend microservices.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         CLIENT                              │
│                     (HTTP/REST/JSON)                        │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                     API GATEWAY :8080                       │
│  ┌────────────────────────────────────────────────────┐    │
│  │  Chi Router                                        │    │
│  │  ├── Global Middleware (Logger, Recoverer, etc.)  │    │
│  │  ├── Public Routes (/auth/*)                      │    │
│  │  └── Protected Routes (/users/*)                  │    │
│  │       └── Auth Middleware (JWT validation)        │    │
│  └────────────────────────────────────────────────────┘    │
│                           │                                 │
│  ┌────────────────────────┴─────────────────────────┐      │
│  │        gRPC Client Connections                   │      │
│  │  ├── AuthServiceClient                           │      │
│  │  └── UserServiceClient                           │      │
│  └──────────────────────────────────────────────────┘      │
└──────────────┬───────────────────────┬────────────────────┘
               │                       │
               ▼                       ▼
    ┌──────────────────┐    ┌──────────────────┐
    │  Auth Service    │    │  User Service    │
    │    :50051        │    │    :50052        │
    │    (gRPC)        │    │    (gRPC)        │
    └──────────────────┘    └──────────────────┘
```

## Key Responsibilities

### 1. **Protocol Translation**
- Accepts HTTP/REST requests from clients
- Translates to gRPC calls for backend services
- Converts gRPC responses back to JSON

### 2. **Authentication & Authorization**
- Validates JWT tokens via Auth Service
- Enforces access control (users can only access their own data)
- Injects user context into requests

### 3. **Cross-Cutting Concerns**
- Request logging
- Panic recovery
- Request ID generation for tracing
- CORS (future)
- Rate limiting (future)

### 4. **Service Orchestration**
- Manages connections to all backend gRPC services
- Handles service failures gracefully
- Provides health check endpoints

## Project Structure

```
api-gateway/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, server setup
├── internal/
│   ├── clients/
│   │   └── grpc_clients.go      # gRPC client connection manager
│   ├── handlers/
│   │   ├── auth_handler.go      # Auth endpoints (register, login)
│   │   └── user_handler.go      # User endpoints (CRUD + addresses)
│   └── middleware/
│       └── auth.go              # JWT validation middleware
├── Dockerfile                   # Multi-stage Docker build
├── go.mod                       # Go module dependencies
└── README.md                    # This file
```

## API Endpoints

### Public Endpoints (No Authentication Required)

| Method | Endpoint               | Description          | Request Body                          |
|--------|------------------------|----------------------|---------------------------------------|
| GET    | `/health`              | Health check         | -                                     |
| POST   | `/api/v1/auth/register`| Register new user    | `{email, password, name}`            |
| POST   | `/api/v1/auth/login`   | Login & get JWT      | `{email, password}`                  |

### Protected Endpoints (JWT Token Required)

| Method | Endpoint                        | Description          | Headers                        |
|--------|---------------------------------|----------------------|--------------------------------|
| GET    | `/api/v1/users/:id`             | Get user profile     | `Authorization: Bearer <token>`|
| PUT    | `/api/v1/users/:id`             | Update user profile  | `Authorization: Bearer <token>`|
| DELETE | `/api/v1/users/:id`             | Delete user          | `Authorization: Bearer <token>`|
| POST   | `/api/v1/users/:id/addresses`   | Add address          | `Authorization: Bearer <token>`|
| GET    | `/api/v1/users/:id/addresses`   | List addresses       | `Authorization: Bearer <token>`|

## How It Works

### Request Flow

#### 1. **Public Route Example: User Login**

```
Client Request:
POST /api/v1/auth/login
Content-Type: application/json
{"email": "user@example.com", "password": "password123"}

↓

API Gateway (auth_handler.go):
├── Parse JSON → LoginRequest struct
├── Validate (email & password not empty)
├── Translate to gRPC:
│   authpb.LoginRequest{Email: "...", Password: "..."}
├── Call: authClient.Login(ctx, grpcReq)
└── Convert gRPC response → JSON

↓

Response:
{
  "token": "eyJhbGc...",
  "user_id": "123",
  "name": "John Doe"
}
```

#### 2. **Protected Route Example: Get User**

```
Client Request:
GET /api/v1/users/123
Authorization: Bearer eyJhbGc...

↓

Auth Middleware (middleware/auth.go):
├── Extract "Bearer eyJhbGc..." from header
├── Call authClient.ValidateToken(token)
├── If valid: Add user_id to context
├── If invalid: Return 401 Unauthorized
└── Pass to next handler

↓

User Handler (handlers/user_handler.go):
├── Get user_id from context
├── Check requested_id == authenticated_id (authorization)
├── Call: userClient.GetUser(ctx, grpcReq)
└── Convert gRPC User → JSON UserResponse

↓

Response:
{
  "id": "123",
  "email": "user@example.com",
  "name": "John Doe",
  "phone": "+1234567890"
}
```

### Authentication Middleware

The auth middleware (`internal/middleware/auth.go`) protects routes by:

1. **Extracting JWT** from `Authorization: Bearer <token>` header
2. **Validating** token format (must start with "Bearer ")
3. **Calling Auth Service** via gRPC to validate the token
4. **Injecting user_id** into request context for handlers
5. **Rejecting** invalid/missing tokens with 401

```go
// Middleware is applied to route groups
r.Route("/users", func(r chi.Router) {
    r.Use(authmw.AuthMiddleware(grpcClients.AuthClient)) // ← All routes below require auth

    r.Get("/{id}", userHandler.GetUser)
    r.Put("/{id}", userHandler.UpdateUser)
    // ... more routes
})
```

## Environment Variables

| Variable           | Description                      | Default              |
|--------------------|----------------------------------|----------------------|
| `AUTH_SERVICE_URL` | Auth service gRPC address        | `localhost:50051`    |
| `USER_SERVICE_URL` | User service gRPC address        | `localhost:50052`    |
| `PORT`             | HTTP port for API Gateway        | `8080`               |

## Running the Gateway

### Local Development

```bash
# From api-gateway directory
go run cmd/server/main.go
```

### With Docker Compose

```bash
# From project root
docker-compose up -d api-gateway
```

The gateway will:
- Connect to auth-service:50051 (gRPC)
- Connect to user-service:50052 (gRPC)
- Listen on :8080 (HTTP)

## Testing the API

### 1. Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "securepass123",
    "name": "Alice Smith"
  }'
```

### 2. Login and Get Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "securepass123"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "abc123",
  "name": "Alice Smith"
}
```

### 3. Access Protected Endpoint

```bash
# Use the token from step 2
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X GET http://localhost:8080/api/v1/users/abc123 \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Test Without Auth (Should Fail)

```bash
curl -X GET http://localhost:8080/api/v1/users/abc123
# Response: 401 Unauthorized - Authorization header required
```

## Key Design Patterns

### 1. **Middleware Pattern**
Higher-order functions that wrap handlers to add functionality:
```go
func AuthMiddleware(client AuthClient) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Validate, then:
            next.ServeHTTP(w, r)
        })
    }
}
```

### 2. **Context Propagation**
User identity flows through the request chain:
```go
// Middleware adds user_id to context
ctx := context.WithValue(r.Context(), userIDKey, "123")

// Handlers extract it
userID := middleware.GetUserID(r.Context())
```

### 3. **Dependency Injection**
Services receive their dependencies via constructors:
```go
authHandler := handlers.NewAuthHandler(grpcClients.AuthClient)
```

### 4. **Graceful Shutdown**
Handles SIGINT/SIGTERM to finish in-flight requests:
```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit // Wait for signal

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
srv.Shutdown(ctx) // Give requests 30s to complete
```

## Error Handling

The gateway returns appropriate HTTP status codes:

| Status Code | Meaning                          | Example                          |
|-------------|----------------------------------|----------------------------------|
| 200 OK      | Successful request               | GET user, login                  |
| 201 Created | Resource created                 | Register, add address            |
| 400 Bad Request | Invalid JSON/missing fields  | Malformed request body           |
| 401 Unauthorized | Missing/invalid token        | No Authorization header          |
| 403 Forbidden | Token valid but no permission   | User A accessing User B's data   |
| 404 Not Found | Resource doesn't exist          | User not found                   |
| 500 Internal Server Error | Backend/gRPC error   | Database down, gRPC call failed  |

## Security Considerations

### ✅ Implemented
- JWT token validation on all protected routes
- Authorization checks (users can only access their own data)
- Custom context keys to prevent value collisions
- Request logging for audit trails

### 🚧 Future Enhancements
- Rate limiting per IP/user
- Request size limits
- CORS configuration
- TLS/HTTPS support
- API key validation for service-to-service
- Request timeouts and circuit breakers

## Monitoring & Observability

### Current Logging
- Structured request logs via Chi's Logger middleware
- gRPC error logging in handlers
- Startup/shutdown logs

### Future Additions
- Prometheus metrics (request count, latency, errors)
- Distributed tracing (OpenTelemetry)
- Health check with service dependency status
- Request/response logging (with PII redaction)

## Troubleshooting

### Gateway can't connect to services

```bash
# Check if services are running
docker-compose ps

# Check gateway logs
docker-compose logs api-gateway

# Verify service URLs
docker-compose exec api-gateway env | grep SERVICE_URL
```

### 401 Unauthorized on all requests

- Verify JWT token is included: `Authorization: Bearer <token>`
- Check token hasn't expired
- Ensure auth-service is running and healthy

### Empty user_id in JWT response

- This happens when login is successful but user-service integration isn't complete
- Auth service creates auth record but can't create user profile
- Verify user-service is running

## Performance Considerations

### gRPC Connection Pooling
- Connections to services are created once at startup
- Reused for all requests (efficient)
- Closed gracefully on shutdown

### HTTP Server Timeouts
```go
ReadTimeout:  15 * time.Second  // Time to read request
WriteTimeout: 15 * time.Second  // Time to write response
IdleTimeout:  60 * time.Second  // Keep-alive timeout
```

### Multi-Stage Docker Build
- Builder stage: 350MB+ (Go toolchain)
- Final image: ~15MB (just binary + Alpine)
- Faster deployments, smaller attack surface

## Future Services Integration

When adding new services (Product, Order, Payment):

1. **Add gRPC client** to `internal/clients/grpc_clients.go`
2. **Create handler** in `internal/handlers/`
3. **Register routes** in `cmd/server/main.go`
4. **Apply middleware** as needed (public vs protected)

Example:
```go
// In main.go
productHandler := handlers.NewProductHandler(grpcClients.ProductClient)

r.Route("/api/v1/products", func(r chi.Router) {
    r.Get("/", productHandler.List)        // Public
    r.Get("/{id}", productHandler.Get)     // Public

    r.Group(func(r chi.Router) {
        r.Use(authmw.AuthMiddleware(grpcClients.AuthClient))
        r.Post("/", productHandler.Create)  // Protected
    })
})
```

---

## Summary

The API Gateway is the **front door** of your microservices architecture. It:
- Simplifies client interaction (one endpoint, REST/JSON)
- Centralizes authentication and cross-cutting concerns
- Translates between HTTP and gRPC protocols
- Enforces security and authorization policies

This design allows backend services to focus on business logic while the gateway handles client-facing concerns.
