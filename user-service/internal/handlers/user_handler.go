package handlers

import (
	"context"

	pb "go-project/proto/user"
	"user-service/internal/service"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserHandler implements the gRPC UserServiceServer interface
type UserHandler struct {
	pb.UnimplementedUserServiceServer // Embedding for forward compatibility
	service                           *service.UserService
}

// NewUserHandler creates a new gRPC handler
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

// CreateUser handles user creation requests
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Call the service layer
	user, err := h.service.CreateUser(req.UserId, req.Email, req.Name, req.Phone)
	if err != nil {
		return &pb.CreateUserResponse{
			User:  nil,
			Error: err.Error(),
		}, nil // Return error in response, not as gRPC error
	}

	// Convert internal model to protobuf
	pbUser := &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Phone:     user.Phone,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Addresses: []*pb.Address{}, // Empty for now, will populate if needed
	}

	return &pb.CreateUserResponse{
		User:  pbUser,
		Error: "",
	}, nil
}

// GetUser handles user retrieval requests
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.service.GetUser(req.UserId)
	if err != nil {
		return &pb.GetUserResponse{
			User:  nil,
			Error: err.Error(),
		}, nil
	}

	// Get user's addresses
	addresses, _ := h.service.GetAddresses(req.UserId)

	// Convert addresses to protobuf
	pbAddresses := make([]*pb.Address, 0, len(addresses))
	for _, addr := range addresses {
		pbAddresses = append(pbAddresses, &pb.Address{
			Id:         addr.ID,
			UserId:     addr.UserID,
			Street:     addr.Street,
			City:       addr.City,
			State:      addr.State,
			PostalCode: addr.PostalCode,
			Country:    addr.Country,
			IsDefault:  addr.IsDefault,
			CreatedAt:  timestamppb.New(addr.CreatedAt),
		})
	}

	pbUser := &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Phone:     user.Phone,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Addresses: pbAddresses,
	}

	return &pb.GetUserResponse{
		User:  pbUser,
		Error: "",
	}, nil
}

// UpdateUser handles user update requests
func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// TODO(human): Implement UpdateUser handler
	// Hints:
	// 1. Call h.service.UpdateUser(req.UserId, req.Name, req.Phone)
	// 2. If error, return pb.UpdateUserResponse with error field set
	// 3. Convert the returned user to pb.User (like in CreateUser)
	// 4. Return pb.UpdateUserResponse with user and empty error
	user, err := h.service.UpdateUser(req.UserId, req.Name, req.Phone)

	if err != nil {
		return &pb.UpdateUserResponse{
			User:  nil,
			Error: err.Error(),
		}, nil
	}

	addresses, _ := h.service.GetAddresses(req.UserId)

	// Convert addresses to protobuf
	pbAddresses := make([]*pb.Address, 0, len(addresses))
	for _, addr := range addresses {
		pbAddresses = append(pbAddresses, &pb.Address{
			Id:         addr.ID,
			UserId:     addr.UserID,
			Street:     addr.Street,
			City:       addr.City,
			State:      addr.State,
			PostalCode: addr.PostalCode,
			Country:    addr.Country,
			IsDefault:  addr.IsDefault,
			CreatedAt:  timestamppb.New(addr.CreatedAt),
		})
	}

	pbUser := &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Phone:     user.Phone,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Addresses: pbAddresses,
	}

	return &pb.UpdateUserResponse{
		User:  pbUser,
		Error: "",
	}, nil
}

// DeleteUser handles user deletion requests
func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// TODO(human): Implement DeleteUser handler
	// Hints:
	// 1. Call h.service.DeleteUser(req.UserId)
	// 2. If error, return pb.DeleteUserResponse{Success: false, Error: err.Error()}
	// 3. If success, return pb.DeleteUserResponse{Success: true, Error: ""}
	err := h.service.DeleteUser(req.UserId)

	if err != nil {
		return &pb.DeleteUserResponse{
			IsDeleted: false,
			Error:     err.Error(),
		}, nil
	}

	return &pb.DeleteUserResponse{
		IsDeleted: true,
		Error:     "",
	}, nil
}

// AddAddress handles address creation requests
func (h *UserHandler) AddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.AddAddressResponse, error) {
	address, err := h.service.AddAddress(
		req.UserId,
		req.Street,
		req.City,
		req.State,
		req.PostalCode,
		req.Country,
		req.IsDefault,
	)

	if err != nil {
		return &pb.AddAddressResponse{
			Address: nil,
			Error:   err.Error(),
		}, nil
	}

	pbAddress := &pb.Address{
		Id:         address.ID,
		UserId:     address.UserID,
		Street:     address.Street,
		City:       address.City,
		State:      address.State,
		PostalCode: address.PostalCode,
		Country:    address.Country,
		IsDefault:  address.IsDefault,
		CreatedAt:  timestamppb.New(address.CreatedAt),
	}

	return &pb.AddAddressResponse{
		Address: pbAddress,
		Error:   "",
	}, nil
}

// GetAddresses handles requests to get all addresses for a user
func (h *UserHandler) GetAddresses(ctx context.Context, req *pb.GetAddressesRequest) (*pb.GetAddressesResponse, error) {
	addresses, err := h.service.GetAddresses(req.UserId)
	if err != nil {
		return &pb.GetAddressesResponse{
			Addresses: nil,
			Error:     err.Error(),
		}, nil
	}

	pbAddresses := make([]*pb.Address, 0, len(addresses))
	for _, addr := range addresses {
		pbAddresses = append(pbAddresses, &pb.Address{
			Id:         addr.ID,
			UserId:     addr.UserID,
			Street:     addr.Street,
			City:       addr.City,
			State:      addr.State,
			PostalCode: addr.PostalCode,
			Country:    addr.Country,
			IsDefault:  addr.IsDefault,
			CreatedAt:  timestamppb.New(addr.CreatedAt),
		})
	}

	return &pb.GetAddressesResponse{
		Addresses: pbAddresses,
		Error:     "",
	}, nil
}
