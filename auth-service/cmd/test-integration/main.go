package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "go-project/proto/auth"
	userpb "go-project/proto/user"
)

func main() {
	// Generate random email to avoid conflicts
	randomID := rand.Intn(100000)
	testEmail := fmt.Sprintf("integrated-test-%d@example.com", randomID)
	testName := fmt.Sprintf("Test User %d", randomID)

	log.Println("🧪 Testing Integrated Auth + User Service Flow")
	log.Printf("📧 Test Email: %s\n", testEmail)

	// Connect to Auth Service
	authConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Auth Service: %v", err)
	}
	defer authConn.Close()
	authClient := authpb.NewAuthServiceClient(authConn)

	// Connect to User Service
	userConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to User Service: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Test 1: Register via Auth Service (should auto-create user profile)
	log.Println("\n📝 Step 1: Registering new user via Auth Service...")
	regResp, err := authClient.Register(ctx, &authpb.RegisterRequest{
		Email:    testEmail,
		Password: "password123",
		Name:     testName,
	})
	if err != nil {
		log.Fatalf("❌ Register failed: %v", err)
	}
	log.Printf("✅ Register successful!")
	log.Printf("   UserID: %s", regResp.UserId)
	log.Printf("   Message: %s", regResp.Message)

	userID := regResp.UserId

	// Test 2: Verify user profile was created in User Service
	log.Println("\n🔍 Step 2: Verifying user profile in User Service...")
	getUserResp, err := userClient.GetUser(ctx, &userpb.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		log.Fatalf("❌ GetUser failed: %v", err)
	}
	if getUserResp.Error != "" {
		log.Fatalf("❌ GetUser error: %s", getUserResp.Error)
	}
	log.Printf("✅ User profile found in User Service!")
	log.Printf("   UserID: %s", getUserResp.User.Id)
	log.Printf("   Email: %s", getUserResp.User.Email)
	log.Printf("   Name: %s", getUserResp.User.Name)

	// Test 3: Login
	log.Println("\n🔐 Step 3: Testing login...")
	loginResp, err := authClient.Login(ctx, &authpb.LoginRequest{
		Email:    testEmail,
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("❌ Login failed: %v", err)
	}
	log.Printf("✅ Login successful!")
	log.Printf("   Token: %s...", loginResp.Token[:30])

	// Test 4: Update user profile
	log.Println("\n✏️  Step 4: Updating user profile...")
	updateResp, err := userClient.UpdateUser(ctx, &userpb.UpdateUserRequest{
		UserId: userID,
		Name:   testName + " (Updated)",
		Phone:  "+1234567890",
	})
	if err != nil {
		log.Fatalf("❌ UpdateUser failed: %v", err)
	}
	if updateResp.Error != "" {
		log.Fatalf("❌ UpdateUser error: %s", updateResp.Error)
	}
	log.Printf("✅ Profile updated!")
	log.Printf("   New Name: %s", updateResp.User.Name)
	log.Printf("   New Phone: %s", updateResp.User.Phone)

	log.Println("\n🎉 All integration tests passed!")
	log.Println("✅ Auth Service and User Service are properly integrated!")
}
