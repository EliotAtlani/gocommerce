package handlers

import (
	"context"
	"errors"
	"log"

	"auth-service/internal/service"
	pb "go-project/proto/auth"
	userpb "go-project/proto/user"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
	userClient  userpb.UserServiceClient // gRPC client for User Service
}

func NewAuthHandler(authService *service.AuthService, userClient userpb.UserServiceClient) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userClient:  userClient,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Step 1: Register user in Auth Service (creates credentials)
	userID, err := h.authService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		return nil, err
	}

	// Step 2: Create user profile in User Service via gRPC

	createUserResp, err := h.userClient.CreateUser(ctx, &userpb.CreateUserRequest{
		UserId: userID,
		Email:  req.Email,
		Name:   req.Name,
		Phone:  "",
	})

	if err != nil {
		return nil, err
	}
	if createUserResp.Error != "" {
		return nil, errors.New(createUserResp.Error)
	}
	log.Printf("✅ User registered: ID=%s, Email=%s", userID, req.Email)

	return &pb.RegisterResponse{
		UserId:  userID,
		Message: "User registered successfully",
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token:  token,
		UserId: "",
		Name:   "",
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse,
	error) {
	userID, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid:  false,
			UserId: "",
			Error:  err.Error(),
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID,
		Error:  "",
	}, nil
}
