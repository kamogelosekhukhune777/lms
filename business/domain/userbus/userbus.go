// Package userbus provides business access to user domain.
package userbus

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"golang.org/x/crypto/bcrypt"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound    = errors.New("user not found")
	ErrUniqueEmail = errors.New("email is not unique")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, usr User) error
}

// Business manages the set of APIs for user access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness constructs a user business API for use.
func NewBusiness(log *logger.Logger, storer Storer) *Business {
	return &Business{
		log:    log,
		storer: storer,
	}
}

// Create adds a new user to the system.
func (b *Business) Create(ctx context.Context, nu NewUser) (User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	usr := User{
		ID:           uuid.New(),
		UserName:     nu.UserName,
		UserEmail:    nu.UserEmail,
		PasswordHash: hash,
		Roles:        nu.Roles,
	}

	if err := b.storer.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}
