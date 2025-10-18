package handlers

import (
	"auth-service/internal/service"
	pb "go-project/proto/auth"
	"context"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	userID, err := h.authService.Register(req.Email, req.Password, req.Name)

	if err != nil {
		return nil, err
	}

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
