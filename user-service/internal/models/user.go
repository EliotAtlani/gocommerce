package models

import "time"

// User represents a user profile in the database
type User struct {
	ID        string
	Email     string
	Name      string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time // Pointer because it can be NULL
}

// Address represents a user's shipping/billing address
type Address struct {
	ID         string
	UserID     string
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
	IsDefault  bool
	CreatedAt  time.Time
}
