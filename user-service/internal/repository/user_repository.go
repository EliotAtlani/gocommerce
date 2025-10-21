package repository

import (
	"database/sql"
	"time"

	"user-service/internal/models"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
// Using an interface allows us to easily mock this for testing
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(userID string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userID string) error

	AddAddress(address *models.Address) error
	GetAddressesByUserID(userID string) ([]*models.Address, error)
}

// PostgresUserRepository implements UserRepository for PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

// CreateUser inserts a new user into the database
func (r *PostgresUserRepository) CreateUser(user *models.User) error {
	// Generate a new UUID for the user if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
		INSERT INTO users (id, email, name, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query, user.ID, user.Email, user.Name, user.Phone, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetUserByID retrieves a user by their ID
func (r *PostgresUserRepository) GetUserByID(userID string) (*models.User, error) {
	user := &models.User{}

	query := `
		SELECT id, email, name, phone, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // User not found
	}

	return user, err
}

// UpdateUser updates an existing user's information
func (r *PostgresUserRepository) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users SET name = $1, phone = $2, updated_at = $3 WHERE id = $4`

	_, err := r.db.Exec(query, user.Name, user.Phone, user.UpdatedAt, user.ID)
	return err
}

// DeleteUser performs a soft delete on a user
func (r *PostgresUserRepository) DeleteUser(userID string) error {
	deleted_at := time.Now()
	query := `
	UPDATE users SET deleted_at = $1 WHERE id = $2`

	_, err := r.db.Exec(query, deleted_at, userID)
	return err
}

// AddAddress adds a new address for a user
func (r *PostgresUserRepository) AddAddress(address *models.Address) error {
	if address.ID == "" {
		address.ID = uuid.New().String()
	}

	address.CreatedAt = time.Now()

	query := `
		INSERT INTO addresses (id, user_id, street, city, state, postal_code, country, is_default, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(query,
		address.ID,
		address.UserID,
		address.Street,
		address.City,
		address.State,
		address.PostalCode,
		address.Country,
		address.IsDefault,
		address.CreatedAt,
	)

	return err
}

// GetAddressesByUserID retrieves all addresses for a user
func (r *PostgresUserRepository) GetAddressesByUserID(userID string) ([]*models.Address, error) {

	query := `
	SELECT id, user_id, street, city, state, postal_code, country, is_default, created_at FROM addresses WHERE user_id = $1
	`

	rows, err := r.db.Query(query, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := []*models.Address{}

	for rows.Next() {
		addr := &models.Address{}
		err := rows.Scan(&addr.ID, &addr.UserID, &addr.Street, &addr.City, &addr.State, &addr.PostalCode, &addr.Country, &addr.IsDefault, &addr.CreatedAt)

		if err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}
