// Package userbus provides business access to user domain.
package userbus

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

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
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
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

func (b *Business) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	user, err := b.storer.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("query: %w", err)
	}

	return user, nil
}

// QueryByEmail finds the user by a specified user email.
func (b *Business) QueryByEmail(ctx context.Context, email mail.Address) (User, error) {
	user, err := b.storer.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	return user, nil
}

// Authenticate finds a user by their email and verifies their passworb. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (b *Business) Authenticate(ctx context.Context, email mail.Address, password string) (User, error) {

	usr, err := b.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return User{}, fmt.Errorf("comparehashandpassword: %w", ErrAuthenticationFailure)
	}

	return usr, nil
}
