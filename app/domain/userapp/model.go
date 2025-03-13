package userapp

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/business/types/name"
	"github.com/kamogelosekhukhune777/lms/business/types/role"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

// NewUser defines the data needed to add a new user.
type NewUser struct {
	Name            string `json:"user_name" validate:"required"`
	Email           string `json:"user_email" validate:"required,email"`
	Role            string `json:"role" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"eqfield=Password"`
}

// Decode implements the decoder interface.
func (app *NewUser) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewUser) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewUser(app NewUser) (userbus.NewUser, error) {

	role, err := role.Parse(app.Role)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	nme, err := name.Parse(app.Name)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	bus := userbus.NewUser{
		UserName:     nme,
		UserEmail:    *addr,
		PasswordHash: app.Password,
	}
	bus.Roles = append(bus.Roles, role) //???

	return bus, nil
}

// =============================================================================

type logInUser struct {
	Email    string `json:"user_email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Decode implements the decoder interface.
func (app *logInUser) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app logInUser) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

//=============================================

type userResponse struct {
	ID           string   `json:"user_id"`
	Name         string   `json:"user_name"`
	Email        string   `json:"user_email"`
	Roles        []string `json:"roles"`
	PasswordHash []byte   `json:"-"`
	DateCreated  string   `json:"CreatedAt"`

	Token string `json:"token"`
}

// Encode implements the web.Encoder interface.
func (ur userResponse) Encode() ([]byte, string, error) {
	b, err := json.Marshal(ur)
	return b, "application/json", err
}

// Headers allows the response to add extra HTTP headers (here, the Authorization header).
func (ur userResponse) Headers() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + ur.Token,
	}
}

// toAppUserWithToken converts the business user and token into a web response.
func toAppUserWithToken(bus *userbus.User, token string) web.Encoder {

	return userResponse{
		ID:           bus.ID.String(),
		Name:         bus.UserName.String(),
		Email:        bus.UserEmail.Address,
		Roles:        role.ParseToString(bus.Roles),
		PasswordHash: bus.PasswordHash,
		DateCreated:  bus.CreatedAt.Format(time.RFC3339),
		Token:        token,
	}
}
