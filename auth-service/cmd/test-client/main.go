package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "go-project/proto/auth"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Test Register
	log.Println("Testing Register...")
	regResp, err := client.Register(ctx, &pb.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	})
	if err != nil {
		log.Fatalf("Register failed: %v", err)
	}
	log.Printf("✅ Register successful: %v\n", regResp)

	// Test Login
	log.Println("\nTesting Login...")
	loginResp, err := client.Login(ctx, &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	log.Printf("✅ Login successful: Token = %s...\n", loginResp.Token[:20])

	// Test ValidateToken
	log.Println("\nTesting ValidateToken...")
	validateResp, err := client.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: loginResp.Token,
	})
	if err != nil {
		log.Fatalf("ValidateToken failed: %v", err)
	}
	log.Printf("✅ ValidateToken successful: Valid=%v, UserID=%s\n", validateResp.Valid, validateResp.UserId)
}
