package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	authpb "go-project/proto/auth"
)

// AuthHandler handles authentication-related HTTP endpoints
type AuthHandler struct {
	authClient authpb.AuthServiceClient
}

func NewAuthHandler(authClient authpb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{authClient: authClient}
}

// Request/Response types for JSON serialization
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type RegisterResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

// TODO(human): Implement Register handler
// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Defer right after the operation it cleans up

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Name == "" || req.Password == "" {
		http.Error(w, "Email,name and password are required", http.StatusBadRequest)
		return
	}

	grpcReq := &authpb.RegisterRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	grpcRes, err := h.authClient.Register(r.Context(), grpcReq)

	if err != nil {
		log.Printf("gRPC Register error: %v", err)
		http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := RegisterResponse{
		UserID:  grpcRes.UserId,
		Message: grpcRes.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// TODO(human): Implement Login handler
// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Your implementation here
	var req LoginRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	grpcReq := &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	grpcRes, err := h.authClient.Login(r.Context(), grpcReq)
	if err != nil {
		log.Printf("gRPC Login error: %v", err)
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized) // 401 for auth failures
		return
	}

	resp := LoginResponse{
		Token:  grpcRes.Token,
		UserID: grpcRes.UserId,
		Name:   grpcRes.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}
