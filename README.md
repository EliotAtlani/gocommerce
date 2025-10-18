# GoCommerce - Microservices Learning Project

A microservices-based e-commerce platform built with Go, gRPC, and PostgreSQL. This project is designed as a hands-on learning experience for mastering Go through real-world microservices architecture.

---

## 🏗️ Architecture

**Monorepo Structure** using Go workspaces (`go.work`):
```
go-project/
├── proto/                    # Shared Protocol Buffer definitions
│   └── auth/                 # Auth service API contract
├── auth-service/             # ✅ Authentication & JWT tokens
├── user-service/             # 🚧 User profile management
├── product-service/          # 🚧 Product catalog
├── order-service/            # 🚧 Order processing
├── payment-service/          # 🚧 Payment handling
├── notification-service/     # 🚧 Email/SMS notifications
├── api-gateway/              # 🚧 REST API gateway
├── docker-compose.yml        # Multi-service orchestration
└── go.work                   # Workspace configuration
```

**Communication:**
- **gRPC** for inter-service communication
- **REST** for external API (via gateway)
- **PostgreSQL** for each service (database-per-service pattern)

---

## 🚀 Tech Stack

| Component | Technology |
|-----------|------------|
| **Language** | Go 1.25+ |
| **Inter-Service** | gRPC + Protocol Buffers |
| **Authentication** | JWT + Bcrypt |
| **Database** | PostgreSQL 15 |
| **Containerization** | Docker + Docker Compose |
| **Message Queue** | Redis / RabbitMQ (planned) |

---

## 📦 Services

### ✅ Auth Service (Complete)
- User registration with bcrypt password hashing
- JWT-based authentication
- Token validation for other services
- **Port:** 50051 (gRPC)
- **Docs:** [auth-service/README.md](auth-service/README.md)

### 🚧 Upcoming Services
- **User Service** - Profile management, preferences
- **Product Service** - Catalog, inventory, search
- **Order Service** - Cart, checkout, order tracking
- **Payment Service** - Payment processing, refunds
- **Notification Service** - Email/SMS notifications
- **API Gateway** - REST endpoints, rate limiting

---

## 🛠️ Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15

### Run All Services

```bash
# Clone the repository
git clone <repo-url>
cd go-project

# Create environment file
cp .env.example .env
# Edit .env with your secrets

# Start all services
docker-compose up -d

# Check logs
docker-compose logs -f auth-service
```

### Test Auth Service

```bash
# From auth-service directory
go run ./cmd/test-client

# Or use grpcurl
grpcurl -plaintext -d '{"email":"test@example.com","password":"pass123","name":"Test"}' \
  localhost:50051 auth.AuthService/Register
```

---

## 📚 Learn More

- **Auth Service Docs:** [auth-service/README.md](auth-service/README.md)
- **Protocol Buffers:** [proto/auth/auth.proto](proto/auth/auth.proto)
- **Go Workspaces:** [go.work](go.work)

---

**Status:** Auth service is production-ready. Other services are under development as part of the learning journey.
