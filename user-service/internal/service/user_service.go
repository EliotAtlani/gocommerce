package service

import (
	"errors"
	"user-service/internal/models"
	"user-service/internal/repository"
)

// UserService contains business logic for user operations
type UserService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates a new user profile
func (s *UserService) CreateUser(userID, email, name, phone string) (*models.User, error) {
	// Validation: email is required
	if email == "" {
		return nil, errors.New("email is required")
	}

	// Validation: name is required
	if name == "" {
		return nil, errors.New("name is required")
	}

	user := &models.User{
		ID:    userID,
		Email: email,
		Name:  name,
		Phone: phone,
	}

	err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(userID string) (*models.User, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(userID, name, phone string) (*models.User, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	user, err := s.repo.GetUserByID(userID)

	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Name = name
	user.Phone = phone

	err_update := s.repo.UpdateUser(user)

	if err_update != nil {
		return nil, err_update
	}

	return user, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	err := s.repo.DeleteUser(userID)

	return err
}

// AddAddress adds a new address for a user
func (s *UserService) AddAddress(userID, street, city, state, postalCode, country string, isDefault bool) (*models.Address, error) {
	// Validation
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if street == "" || city == "" || postalCode == "" || country == "" {
		return nil, errors.New("street, city, postal code, and country are required")
	}

	address := &models.Address{
		UserID:     userID,
		Street:     street,
		City:       city,
		State:      state,
		PostalCode: postalCode,
		Country:    country,
		IsDefault:  isDefault,
	}

	err := s.repo.AddAddress(address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

// GetAddresses retrieves all addresses for a user
func (s *UserService) GetAddresses(userID string) ([]*models.Address, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return s.repo.GetAddressesByUserID(userID)
}
