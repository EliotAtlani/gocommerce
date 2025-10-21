# **GoCommerce - Distributed E-Commerce Microservices Platform**

## **Project Overview**

Build a production-ready e-commerce backend using Go microservices that communicate via gRPC and REST APIs. This project demonstrates real-world distributed systems architecture with authentication, asynchronous processing, and inter-service communication.

## **System Architecture**

### **Core Microservices**

1. **API Gateway Service** (REST)

   - Single entry point for all client requests
   - Routes requests to appropriate microservices
   - Rate limiting and request validation
   - HTTP to gRPC translation layer

2. **Authentication Service** (gRPC + REST)

   - User registration and login
   - JWT token generation and validation
   - Password hashing (bcrypt)
   - Token refresh mechanism
   - Role-based access control (RBAC)

3. **User Service** (gRPC)

   - User profile management
   - User preferences and settings
   - Address management
   - PostgreSQL database

4. **Product Service** (gRPC)

   - Product catalog management
   - Inventory tracking
   - Product search and filtering
   - Category management
   - PostgreSQL database

5. **Order Service** (gRPC)

   - Order creation and management
   - Order status tracking
   - Order history
   - Integration with payment and inventory
   - PostgreSQL database

6. **Payment Service** (gRPC)

   - Payment processing simulation
   - Payment status tracking
   - Refund handling
   - Integration with task queue for async processing

7. **Notification Service** (gRPC + Task Queue Consumer)
   - Email notifications (order confirmations, shipping updates)
   - SMS notifications (optional)
   - Push notifications
   - Consumes messages from task queue

## **Technical Stack**

### **Core Technologies**

- **Language**: Go 1.21+
- **Communication**: gRPC with Protocol Buffers
- **REST API**: Chi/Gin/Fiber router
- **Authentication**: JWT tokens
- **Databases**: PostgreSQL (per service)
- **Task Queue**: RabbitMQ or Redis Queue
- **Caching**: Redis (optional)

### **Go Libraries**

- `google.golang.org/grpc` - gRPC framework
- `google.golang.org/protobuf` - Protocol Buffers
- `github.com/golang-jwt/jwt` - JWT authentication
- `github.com/lib/pq` or `gorm.io/gorm` - PostgreSQL driver/ORM
- `github.com/rabbitmq/amqp091-go` - RabbitMQ client
- `golang.org/x/crypto/bcrypt` - Password hashing

## **Key Features to Implement**

### **Phase 1: Foundation**

- ✅ Project structure setup (multi-module workspace)
- ✅ Protocol Buffer definitions for all services
- ✅ Database schemas and migrations
- ✅ Authentication service with JWT
- ✅ API Gateway with basic routing

### **Phase 2: Core Services**

- ✅ User service with CRUD operations
- ✅ Product service with inventory management
- ✅ Service-to-service gRPC communication
- ✅ Middleware for authentication in API Gateway
- ✅ Error handling and logging

### **Phase 3: Business Logic**

- ✅ Order creation workflow (multi-service transaction)
- ✅ Inventory deduction on order placement
- ✅ Payment processing integration
- ✅ Order status state machine

### **Phase 4: Asynchronous Processing**

- ✅ RabbitMQ/Redis Queue setup
- ✅ Task producers in services (Order, Payment)
- ✅ Notification service as consumer
- ✅ Email/SMS sending after order events
- ✅ Retry mechanisms and dead-letter queues

### **Phase 5: Production Readiness**

- ✅ Docker containerization for each service
- ✅ Docker Compose for local development
- ✅ Health check endpoints
- ✅ Graceful shutdown handling
- ✅ Configuration management (environment variables)
- ✅ Structured logging (zerolog/zap)
- ✅ Basic metrics and monitoring setup

## **Learning Outcomes**

By completing this project, you will learn:

1. **Go Fundamentals**

   - Concurrency patterns (goroutines, channels)
   - Context management
   - Error handling best practices
   - Package organization

2. **Microservices Architecture**

   - Service decomposition principles
   - Inter-service communication patterns
   - Data consistency in distributed systems
   - Service discovery concepts

3. **gRPC & Protocol Buffers**

   - `.proto` file definitions
   - Code generation
   - Unary and streaming RPCs
   - Error handling in gRPC

4. **Authentication & Security**

   - JWT implementation
   - Secure password storage
   - API authentication middleware
   - Authorization patterns

5. **Asynchronous Processing**

   - Message queue patterns
   - Producer-consumer architecture
   - Reliable message delivery
   - Background job processing

6. **Database Management**

   - Database per service pattern
   - Migrations
   - Connection pooling
   - Transaction handling

7. **DevOps Basics**
   - Containerization
   - Multi-container orchestration
   - Environment configuration
   - Service health monitoring

## **Project Structure**

```
go-commerce/
├── api-gateway/          # REST API Gateway
├── auth-service/         # Authentication Service
├── user-service/         # User Management
├── product-service/      # Product Catalog
├── order-service/        # Order Management
├── payment-service/      # Payment Processing
├── notification-service/ # Notifications (Queue Consumer)
├── proto/               # Shared Protocol Buffer definitions
├── pkg/                 # Shared packages (middleware, utils)
├── docker-compose.yml   # Local development setup
└── README.md
```

## **Success Criteria**

- All services run independently in Docker containers
- API Gateway successfully routes REST requests to gRPC services
- JWT authentication works across all protected endpoints
- Complete order flow: create user → browse products → place order → process payment → send notification
- Task queue successfully processes asynchronous jobs
- Services handle failures gracefully with proper error responses

---

This project provides a comprehensive learning path through Go while building something practical and portfolio-worthy.
