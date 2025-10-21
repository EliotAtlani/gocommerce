package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "go-project/proto/user"
)

func main() {
	// Connect to User Service on port 50052
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Test 1: CreateUser
	log.Println("📝 Testing CreateUser...")
	createResp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		UserId: "user-12345",
		Email:  "john.doe@example.com",
		Name:   "John Doe",
		Phone:  "+1234567890",
	})
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	if createResp.Error != "" {
		log.Fatalf("CreateUser error: %s", createResp.Error)
	}
	log.Printf("✅ CreateUser successful: UserID=%s, Email=%s\n", createResp.User.Id, createResp.User.Email)

	// Test 2: GetUser
	log.Println("\n🔍 Testing GetUser...")
	getResp, err := client.GetUser(ctx, &pb.GetUserRequest{
		UserId: "user-12345",
	})
	if err != nil {
		log.Fatalf("GetUser failed: %v", err)
	}
	if getResp.Error != "" {
		log.Fatalf("GetUser error: %s", getResp.Error)
	}
	log.Printf("✅ GetUser successful: Name=%s, Phone=%s\n", getResp.User.Name, getResp.User.Phone)

	// Test 3: UpdateUser
	log.Println("\n✏️  Testing UpdateUser...")
	updateResp, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
		UserId: "user-12345",
		Name:   "John Updated Doe",
		Phone:  "+9876543210",
	})
	if err != nil {
		log.Fatalf("UpdateUser failed: %v", err)
	}
	if updateResp.Error != "" {
		log.Fatalf("UpdateUser error: %s", updateResp.Error)
	}
	log.Printf("✅ UpdateUser successful: NewName=%s, NewPhone=%s\n", updateResp.User.Name, updateResp.User.Phone)

	// Test 4: AddAddress
	log.Println("\n🏠 Testing AddAddress...")
	addrResp, err := client.AddAddress(ctx, &pb.AddAddressRequest{
		UserId:     "user-12345",
		Street:     "123 Main St",
		City:       "San Francisco",
		State:      "CA",
		PostalCode: "94102",
		Country:    "USA",
		IsDefault:  true,
	})
	if err != nil {
		log.Fatalf("AddAddress failed: %v", err)
	}
	if addrResp.Error != "" {
		log.Fatalf("AddAddress error: %s", addrResp.Error)
	}
	log.Printf("✅ AddAddress successful: %s, %s, %s %s\n",
		addrResp.Address.Street, addrResp.Address.City, addrResp.Address.State, addrResp.Address.PostalCode)

	// Add another address
	log.Println("\n🏢 Testing AddAddress (second address)...")
	addr2Resp, err := client.AddAddress(ctx, &pb.AddAddressRequest{
		UserId:     "user-12345",
		Street:     "456 Work Ave",
		City:       "Palo Alto",
		State:      "CA",
		PostalCode: "94301",
		Country:    "USA",
		IsDefault:  false,
	})
	if err != nil {
		log.Fatalf("AddAddress (2nd) failed: %v", err)
	}
	if addr2Resp.Error != "" {
		log.Fatalf("AddAddress (2nd) error: %s", addr2Resp.Error)
	}
	log.Printf("✅ AddAddress (2nd) successful: %s, %s\n", addr2Resp.Address.Street, addr2Resp.Address.City)

	// Test 5: GetAddresses
	log.Println("\n📍 Testing GetAddresses...")
	addrsResp, err := client.GetAddresses(ctx, &pb.GetAddressesRequest{
		UserId: "user-12345",
	})
	if err != nil {
		log.Fatalf("GetAddresses failed: %v", err)
	}
	if addrsResp.Error != "" {
		log.Fatalf("GetAddresses error: %s", addrsResp.Error)
	}
	log.Printf("✅ GetAddresses successful: Found %d addresses\n", len(addrsResp.Addresses))
	for i, addr := range addrsResp.Addresses {
		log.Printf("   Address %d: %s, %s (Default: %v)\n", i+1, addr.Street, addr.City, addr.IsDefault)
	}

	// Test 6: GetUser with addresses
	log.Println("\n👤 Testing GetUser (with addresses)...")
	getUserResp, err := client.GetUser(ctx, &pb.GetUserRequest{
		UserId: "user-12345",
	})
	if err != nil {
		log.Fatalf("GetUser failed: %v", err)
	}
	if getUserResp.Error != "" {
		log.Fatalf("GetUser error: %s", getUserResp.Error)
	}
	log.Printf("✅ GetUser with addresses: %s has %d addresses\n",
		getUserResp.User.Name, len(getUserResp.User.Addresses))

	// Test 7: DeleteUser (soft delete)
	log.Println("\n🗑️  Testing DeleteUser...")
	deleteResp, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{
		UserId: "user-12345",
	})
	if err != nil {
		log.Fatalf("DeleteUser failed: %v", err)
	}
	if deleteResp.Error != "" {
		log.Fatalf("DeleteUser error: %s", deleteResp.Error)
	}
	log.Printf("✅ DeleteUser successful: Deleted=%v\n", deleteResp.IsDeleted)

	// Test 8: Try to get deleted user (should fail)
	log.Println("\n❌ Testing GetUser after delete (should fail)...")
	getDeletedResp, err := client.GetUser(ctx, &pb.GetUserRequest{
		UserId: "user-12345",
	})
	if err != nil {
		log.Fatalf("GetUser failed: %v", err)
	}
	if getDeletedResp.Error != "" {
		log.Printf("✅ Expected error received: %s\n", getDeletedResp.Error)
	} else {
		log.Println("⚠️  Warning: Deleted user was still retrieved!")
	}

	log.Println("\n🎉 All tests completed successfully!")
}
