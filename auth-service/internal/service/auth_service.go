package service

import (
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account
func (s *AuthService) Register(email, password, name string) (string, error) {
	// TODO(human): Implement registration logic
	existingUser, err := s.repo.GetUserByEmail(email)

	if err != nil {
		return "", err
	}
	if existingUser != nil {
		return "", errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	newUser := &models.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	err = s.repo.CreateUser(newUser)

	if err != nil {
		return "", err
	}

	return newUser.ID, nil

}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(email, password string) (string, error) {
	// TODO(human): Implement login logic
	user, err := s.repo.GetUserByEmail(email)

	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	// verify pwd
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

// ValidateToken checks if a JWT token is valid and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return "", errors.New("Invalid token")
	}

	userID := claims["user_id"]
	userIDStr, ok := userID.(string)

	if !ok {
		return "", errors.New("Invalid token claims")
	}

	return userIDStr, nil
}
