package userbus

import (
	"net/mail"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/types/name"
	"github.com/kamogelosekhukhune777/lms/business/types/role"
)

// User represents information about an individual user.
type User struct {
	ID           uuid.UUID
	UserName     name.Name
	UserEmail    mail.Address
	PasswordHash []byte
	Roles        []role.Role
	Enabled      bool
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	UserName  name.Name
	UserEmail mail.Address
	Password  string
	Roles     []role.Role
}

type LogInUser struct {
	UserEmail mail.Address
	Password  string
}
