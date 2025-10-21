package clients

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "go-project/proto/auth"
	userpb "go-project/proto/user"
)

// GRPCClients holds connections to all backend microservices
// This struct centralizes our gRPC client management, making it easy to
// inject into handlers and manage connection lifecycle
type GRPCClients struct {
	AuthClient authpb.AuthServiceClient
	UserClient userpb.UserServiceClient
	// Future: ProductClient, OrderClient, etc.

	// Keep connection references for cleanup
	authConn *grpc.ClientConn
	userConn *grpc.ClientConn
}

// NewGRPCClients establishes connections to all backend services
// In production, these addresses would come from service discovery (Consul, Kubernetes DNS, etc.)
func NewGRPCClients(authServiceAddr, userServiceAddr string) (*GRPCClients, error) {
	// Connect to Auth Service
	authConn, err := grpc.Dial(
		authServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// TODO: Add interceptors for logging, tracing, retry logic
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	// Connect to User Service
	userConn, err := grpc.Dial(
		userServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		authConn.Close() // Clean up first connection
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	log.Printf("✓ Connected to Auth Service at %s", authServiceAddr)
	log.Printf("✓ Connected to User Service at %s", userServiceAddr)

	return &GRPCClients{
		AuthClient: authpb.NewAuthServiceClient(authConn),
		UserClient: userpb.NewUserServiceClient(userConn),
		authConn:   authConn,
		userConn:   userConn,
	}, nil
}

// Close gracefully closes all gRPC connections
// This should be called when the gateway shuts down
func (c *GRPCClients) Close() {
	if c.authConn != nil {
		if err := c.authConn.Close(); err != nil {
			log.Printf("Error closing auth connection: %v", err)
		}
	}
	if c.userConn != nil {
		if err := c.userConn.Close(); err != nil {
			log.Printf("Error closing user connection: %v", err)
		}
	}
	log.Println("All gRPC connections closed")
}
