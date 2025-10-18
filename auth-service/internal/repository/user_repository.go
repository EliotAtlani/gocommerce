package repository

import (
	"auth-service/internal/models"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	query := `INSERT INTO users (id, email, password, name, created_at, last_login) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query, user.ID, user.Email, user.Password, user.Name, user.CreatedAt, user.LastLogin)

	return err
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	query := `SELECT id, email, password, name, created_at, last_login FROM users WHERE email  = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.LastLogin,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (r *PostgresUserRepository) GetUserByID(id string) (*models.User, error) {
	user := &models.User{}

	query := `SELECT id, email, password, name, created_at, last_login FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.LastLogin,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}
