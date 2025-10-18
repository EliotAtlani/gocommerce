# Auth Service

Microservice for user authentication and authorization using gRPC, PostgreSQL, and JWT tokens.

---

## Overview

The **Auth Service** handles:
- ✅ User registration with bcrypt password hashing
- ✅ User login with JWT token generation
- ✅ Token validation for other microservices

**Tech Stack:**
- **gRPC** - Inter-service communication
- **Protocol Buffers** - API contracts
- **PostgreSQL** - User data storage
- **Bcrypt** - Password hashing
- **JWT** - Stateless authentication tokens
- **Docker** - Containerization

---

## Architecture

### Layered Design

```
┌─────────────────────────────────────┐
│  gRPC Client (Other Services)       │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│  Handlers (Presentation Layer)      │
│  • auth_handler.go                  │
│  • Converts protobuf ↔ service      │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│  Service (Business Logic)           │
│  • auth_service.go                  │
│  • Password hashing (bcrypt)        │
│  • JWT generation & validation      │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│  Repository (Data Access)           │
│  • user_repository.go               │
│  • SQL queries (CRUD operations)    │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│  PostgreSQL Database                │
│  • users table                      │
└─────────────────────────────────────┘
```

### Directory Structure

```
auth-service/
├── cmd/
│   ├── server/           # Main server entry point
│   │   └── main.go
│   └── test-client/      # Test client for manual testing
│       └── main.go
├── internal/
│   ├── handlers/         # gRPC handlers
│   │   └── auth_handler.go
│   ├── service/          # Business logic
│   │   └── auth_service.go
│   ├── repository/       # Database layer
│   │   └── user_repository.go
│   └── models/           # Data structures
│       └── user.go
├── migrations/           # SQL schema migrations
│   └── 001_create_users_table.sql
├── Dockerfile            # Multi-stage Docker build
├── go.mod                # Go module dependencies
└── README.md             # This file
```

---

## gRPC API

### Service Definition

The service contract is defined in `../proto/auth/auth.proto`:

```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
}
```

### Available Methods

#### 1. Register
Creates a new user account.

**Request:**
```protobuf
message RegisterRequest {
  string email = 1;
  string password = 2;
  string name = 3;
}
```

**Response:**
```protobuf
message RegisterResponse {
  string user_id = 1;
  string message = 2;
}
```

**Example:**
```bash
grpcurl -plaintext -d '{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}' localhost:50051 auth.AuthService/Register
```

---

#### 2. Login
Authenticates a user and returns a JWT token.

**Request:**
```protobuf
message LoginRequest {
  string email = 1;
  string password = 2;
}
```

**Response:**
```protobuf
message LoginResponse {
  string token = 1;       # JWT access token
  string user_id = 2;
  string name = 3;
}
```

**Example:**
```bash
grpcurl -plaintext -d '{
  "email": "user@example.com",
  "password": "password123"
}' localhost:50051 auth.AuthService/Login
```

---

#### 3. ValidateToken
Validates a JWT token (used by other microservices).

**Request:**
```protobuf
message ValidateTokenRequest {
  string token = 1;
}
```

**Response:**
```protobuf
message ValidateTokenResponse {
  bool valid = 1;
  string user_id = 2;
  string error = 3;
}
```

**Example:**
```bash
grpcurl -plaintext -d '{
  "token": "eyJhbGc..."
}' localhost:50051 auth.AuthService/ValidateToken
```

---

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

**Fields:**
- `id` - UUID generated on user creation
- `email` - Unique user email (used for login)
- `password` - Bcrypt hashed password (never stored in plain text)
- `name` - User's display name
- `created_at` - Account creation timestamp
- `last_login` - Last login timestamp

---

## Running the Service

### Prerequisites

- **Go 1.21+**
- **PostgreSQL 15**
- **Docker & Docker Compose** (optional)

### Local Development

#### 1. Set Environment Variables

```bash
export DATABASE_URL="postgres://authuser:authpass@localhost:5432/authdb?sslmode=disable"
export JWT_SECRET="your-secret-key-here"
```

#### 2. Run Database Migrations

```bash
psql -h localhost -U authuser -d authdb -f migrations/001_create_users_table.sql
```

#### 3. Run the Server

```bash
go run ./cmd/server
```

**Output:**
```
2025/10/18 00:00:00 Auth service listening on :50051
```

---

### Docker Compose (Recommended)

#### 1. Start All Services

```bash
# From project root
docker-compose up -d
```

This starts:
- ✅ PostgreSQL (port 5432)
- ✅ Auth Service (port 50051)

#### 2. Check Logs

```bash
docker-compose logs -f auth-service
```

#### 3. Stop Services

```bash
docker-compose down
```

---

## Testing

### Manual Testing with Test Client

```bash
go run ./cmd/test-client
```

**Output:**
```
2025/10/18 00:00:00 Testing Register...
2025/10/18 00:00:00 ✅ Register successful: user_id:"abc-123" message:"User registered successfully"

2025/10/18 00:00:00 Testing Login...
2025/10/18 00:00:00 ✅ Login successful: Token = eyJhbGciOiJIUzI1NiIs...

2025/10/18 00:00:00 Testing ValidateToken...
2025/10/18 00:00:00 ✅ ValidateToken successful: Valid=true, UserID=abc-123
```

### Using grpcurl

Install grpcurl:
```bash
brew install grpcurl
```

Test the service:
```bash
# List available services
grpcurl -plaintext localhost:50051 list

# Register a user
grpcurl -plaintext -d '{"email":"test@example.com","password":"pass123","name":"Test"}' \
  localhost:50051 auth.AuthService/Register

# Login
grpcurl -plaintext -d '{"email":"test@example.com","password":"pass123"}' \
  localhost:50051 auth.AuthService/Login
```

---

## How Other Services Use This

### Example: Order Service Validating Tokens

```go
package clients

import (
    "context"
    "google.golang.org/grpc"
    pb "go-project/proto/auth"  // Shared proto
)

type AuthClient struct {
    client pb.AuthServiceClient
}

func NewAuthClient(address string) (*AuthClient, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }

    return &AuthClient{
        client: pb.NewAuthServiceClient(conn),
    }, nil
}

func (ac *AuthClient) ValidateToken(token string) (string, bool, error) {
    resp, err := ac.client.ValidateToken(context.Background(),
        &pb.ValidateTokenRequest{Token: token})

    if err != nil {
        return "", false, err
    }

    return resp.UserId, resp.Valid, nil
}
```

---

## Security Features

✅ **Bcrypt Password Hashing**
- Passwords are hashed with bcrypt (cost factor: 10)
- Plain text passwords are never stored
- Salting is automatic

✅ **JWT Tokens**
- Signed with HS256 algorithm
- 24-hour expiration
- Contains user_id claim

✅ **Input Validation**
- Email uniqueness enforced by database constraint
- Generic error messages (prevents account enumeration)

✅ **SQL Injection Protection**
- Parameterized queries ($1, $2 placeholders)

---

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@host:5432/dbname?sslmode=disable` |
| `JWT_SECRET` | Secret key for signing JWT tokens | Generated with `openssl rand -base64 32` |

---

## Dependencies

```go
require (
    github.com/golang-jwt/jwt/v5       // JWT token generation
    github.com/google/uuid             // UUID generation
    github.com/lib/pq                  // PostgreSQL driver
    golang.org/x/crypto/bcrypt         // Password hashing
    google.golang.org/grpc             // gRPC framework
    google.golang.org/protobuf         // Protocol Buffers
)
```

---

## Troubleshooting

### Service won't start

**Check database connection:**
```bash
docker exec -it auth-postgres psql -U authuser -d authdb -c "SELECT 1;"
```

**Check if port 50051 is in use:**
```bash
lsof -i :50051
```

### User registration fails

**Check if email already exists:**
```bash
docker exec -it auth-postgres psql -U authuser -d authdb \
  -c "SELECT email FROM users WHERE email = 'test@example.com';"
```

### Token validation always fails

**Check JWT_SECRET matches:**
- Ensure the same secret is used for signing and validating
- Tokens signed with one secret can't be validated with another

---

## Next Steps

- [ ] Add password reset functionality
- [ ] Implement refresh tokens
- [ ] Add rate limiting
- [ ] Add email verification
- [ ] Implement OAuth2 providers
- [ ] Add observability (metrics, tracing)

---

## Learn More

- [Protocol Buffers](https://protobuf.dev/)
- [gRPC in Go](https://grpc.io/docs/languages/go/)
- [JWT Best Practices](https://datatracker.ietf.org/doc/html/rfc8725)
- [Go Database/SQL Tutorial](https://go.dev/doc/database/sql-injection)
