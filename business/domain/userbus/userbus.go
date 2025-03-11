// Package userbus provides business access to user domain.
package userbus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"golang.org/x/crypto/bcrypt"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, usr User) error
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
}

// Business manages the set of APIs for user access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// Create adds a new user to the system.
func (b *Business) Create(ctx context.Context, nu NewUser) (User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	usr := User{
		ID:           uuid.New(),
		UserName:     nu.UserName,
		UserEmail:    nu.UserEmail,
		PasswordHash: hash,
		Roles:        nu.Roles,
		CreatedAt:    now,
	}

	if err := b.storer.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}

// QueryByID finds the user by the specified ID.
func (b *Business) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	user, err := b.storer.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("query: userID[%s]: %w", userID, err)
	}

	return user, nil
}
