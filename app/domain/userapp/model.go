package userapp

import (
	"fmt"
	"net/mail"

	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/business/types/name"
	"github.com/kamogelosekhukhune777/lms/business/types/role"
)

// User represents information about an individual user.
type User struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
	PasswordHash []byte   `json:"-"`
	Enabled      bool     `json:"enabled"`
}

func toAppUser(bus userbus.User) User {
	return User{
		ID:           bus.ID.String(),
		Name:         bus.UserName.String(),
		Email:        bus.UserEmail.Address,
		Roles:        role.ParseToString(bus.Roles),
		PasswordHash: bus.PasswordHash,
		Enabled:      bus.Enabled,
	}
}

// =============================================================================

// NewUser defines the data needed to add a new user.
type NewUser struct {
	UserName  string   `json:"name" validate:"required"`
	UserEmail string   `json:"email" validate:"required,email"`
	Roles     []string `json:"roles" validate:"required"`
	Password  string   `json:"password" validate:"required"`
	//PasswordConfirm string   `json:"passwordConfirm" validate:"eqfield=Password"`
}

// Validate checks the data in the model is considered clean.

func (app NewUser) Validate() error {
	if err := errs.Check(app); err != nil {
		return errs.Newf(errs.InvalidArgument, "validate: %s", err)
	}

	return nil
}

func toBusNewUser(app NewUser) (userbus.NewUser, error) {
	roles, err := role.ParseMany(app.Roles)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	addr, err := mail.ParseAddress(app.UserEmail)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	nme, err := name.Parse(app.UserName)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	bus := userbus.NewUser{
		UserName:  nme,
		UserEmail: *addr,
		Roles:     roles,
		Password:  app.Password,
	}

	return bus, nil
}

// =============================================================================

// logInUser defines the data needed to logInUser a user.
type logInUser struct {
	UserEmail string `json:"user_email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

// Validate checks the data in the model is considered clean.

func (app logInUser) Validate() error {
	if err := errs.Check(app); err != nil {
		return errs.Newf(errs.InvalidArgument, "validate: %s", err)
	}

	return nil
}

func tologInUser(app logInUser) (userbus.LogInUser, error) {
	addr, err := mail.ParseAddress(app.UserEmail)
	if err != nil {
		return userbus.LogInUser{}, fmt.Errorf("parse: %w", err)
	}

	bus := userbus.LogInUser{
		UserEmail: *addr,
		Password:  app.Password,
	}

	return bus, nil
}
