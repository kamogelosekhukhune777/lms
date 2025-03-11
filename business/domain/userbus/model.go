package userbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/types/role"
)

// User represents information about an individual user.
type User struct {
	ID           uuid.UUID
	UserName     string
	UserEmail    string
	PasswordHash []byte
	Roles        []role.Role
	CreatedAt    time.Time
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	UserName     string
	UserEmail    string
	PasswordHash string
	Roles        []role.Role
}
