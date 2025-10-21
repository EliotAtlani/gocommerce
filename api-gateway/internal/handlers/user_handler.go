package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"api-gateway/internal/middleware"
	userpb "go-project/proto/user"
)

// UserHandler handles user-related HTTP endpoints
type UserHandler struct {
	userClient userpb.UserServiceClient
}

func NewUserHandler(userClient userpb.UserServiceClient) *UserHandler {
	return &UserHandler{userClient: userClient}
}

// Request/Response types
type UpdateUserRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type AddAddressRequest struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	IsDefault  bool   `json:"is_default"`
}

type AddressResponse struct {
	ID         string `json:"id"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	IsDefault  bool   `json:"is_default"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

// Helper function

func authorizeUserAccess(r *http.Request) (string, error) {
	requestedID := chi.URLParam(r, "id")
	authenticatedID := middleware.GetUserID(r.Context())

	if authenticatedID == "" {
		return "", errors.New("no authenticated user found")
	}

	if requestedID != authenticatedID {
		return "", errors.New("forbidden: cannot access other users' data")
	}

	return requestedID, nil
}

// GetUser handles GET /api/v1/users/:id
// This is a protected route - user_id will be in the context from auth middleware
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	requestedUserID, err := authorizeUserAccess(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Call the user service via gRPC
	grpcResp, err := h.userClient.GetUser(r.Context(), &userpb.GetUserRequest{
		UserId: requestedUserID,
	})
	if err != nil {
		log.Printf("gRPC GetUser error: %v", err)
		http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if grpcResp.Error != "" {
		http.Error(w, grpcResp.Error, http.StatusNotFound)
		return
	}

	// Convert protobuf User to JSON response
	resp := UserResponse{
		ID:    grpcResp.User.Id,
		Email: grpcResp.User.Email,
		Name:  grpcResp.User.Name,
		Phone: grpcResp.User.Phone,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UpdateUser handles PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	requestedUserID, err := authorizeUserAccess(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var req UpdateUserRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Name == "" && req.Phone == "" {
		http.Error(w, "At least one field (name or phone) must be provided", http.StatusBadRequest)
		return
	}

	grpcResp, err := h.userClient.UpdateUser(r.Context(), &userpb.UpdateUserRequest{
		UserId: requestedUserID,
		Name:   req.Name,
		Phone:  req.Phone,
	})

	if err != nil {
		log.Printf("gRPC UpdateUser error: %v", err)
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if grpcResp.Error != "" {
		http.Error(w, grpcResp.Error, http.StatusNotFound)
		return
	}

	resp := UserResponse{
		ID:    grpcResp.User.Id,
		Email: grpcResp.User.Email,
		Name:  grpcResp.User.Name,
		Phone: grpcResp.User.Phone,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DeleteUser handles DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Your implementation here
	requestedUserID, err := authorizeUserAccess(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Call the user service via gRPC
	grpcResp, err := h.userClient.DeleteUser(r.Context(), &userpb.DeleteUserRequest{
		UserId: requestedUserID,
	})

	if err != nil {
		log.Printf("gRPC DeleteUser error: %v", err)
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if grpcResp.Error != "" {
		http.Error(w, grpcResp.Error, http.StatusNotFound)
		return
	}

	// Convert protobuf User to JSON response
	resp := DeleteUserResponse{
		Message: "User deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// AddAddress handles POST /api/v1/users/:id/addresses
func (h *UserHandler) AddAddress(w http.ResponseWriter, r *http.Request) {
	// Authorize: user can only add addresses to their own profile
	requestedUserID, err := authorizeUserAccess(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Parse the address data from request body
	var req AddAddressRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate required fields for an address
	if req.Street == "" || req.City == "" || req.Country == "" {
		http.Error(w, "Street, city, and country are required", http.StatusBadRequest)
		return
	}

	// Call gRPC service to add the address
	grpcResp, err := h.userClient.AddAddress(r.Context(), &userpb.AddAddressRequest{
		UserId:     requestedUserID,
		Street:     req.Street,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		Country:    req.Country,
		IsDefault:  req.IsDefault,
	})

	if err != nil {
		log.Printf("gRPC AddAddress error: %v", err)
		http.Error(w, "Failed to add address: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if grpcResp.Error != "" {
		http.Error(w, grpcResp.Error, http.StatusBadRequest)
		return
	}

	// Convert protobuf Address to JSON response
	resp := AddressResponse{
		ID:         grpcResp.Address.Id,
		Street:     grpcResp.Address.Street,
		City:       grpcResp.Address.City,
		State:      grpcResp.Address.State,
		PostalCode: grpcResp.Address.PostalCode,
		Country:    grpcResp.Address.Country,
		IsDefault:  grpcResp.Address.IsDefault,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created for new resource
	json.NewEncoder(w).Encode(resp)
}

// GetAddresses handles GET /api/v1/users/:id/addresses
func (h *UserHandler) GetAddresses(w http.ResponseWriter, r *http.Request) {
	// Authorize: user can only view their own addresses
	requestedUserID, err := authorizeUserAccess(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Call gRPC service to get all addresses for this user
	grpcResp, err := h.userClient.GetAddresses(r.Context(), &userpb.GetAddressesRequest{
		UserId: requestedUserID,
	})

	if err != nil {
		log.Printf("gRPC GetAddresses error: %v", err)
		http.Error(w, "Failed to get addresses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if grpcResp.Error != "" {
		http.Error(w, grpcResp.Error, http.StatusNotFound)
		return
	}

	// Convert repeated protobuf Addresses to JSON array
	// Note: We need to handle the case where addresses might be empty
	addresses := make([]AddressResponse, 0, len(grpcResp.Addresses))
	for _, addr := range grpcResp.Addresses {
		addresses = append(addresses, AddressResponse{
			ID:         addr.Id,
			Street:     addr.Street,
			City:       addr.City,
			State:      addr.State,
			PostalCode: addr.PostalCode,
			Country:    addr.Country,
			IsDefault:  addr.IsDefault,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addresses)
}
