package userdb

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/sqldb/dbarray"
	"github.com/kamogelosekhukhune777/lms/business/types/name"
	"github.com/kamogelosekhukhune777/lms/business/types/role"
)

type user struct {
	ID           uuid.UUID      `db:"user_id"`
	Name         string         `db:"user_name"`
	Email        string         `db:"user_email"`
	Roles        dbarray.String `db:"roles"`
	PasswordHash []byte         `db:"password_hash"`
	CreatedAt    time.Time      `db:"created_at"`
}

func toDBUser(bus userbus.User) user {
	return user{
		ID:           bus.ID,
		Name:         bus.UserName.String(),
		Email:        bus.UserEmail.Address,
		Roles:        role.ParseToString(bus.Roles),
		PasswordHash: bus.PasswordHash,
		CreatedAt:    bus.CreatedAt.UTC(),
	}
}

func toBusUser(db user) (userbus.User, error) {
	addr := mail.Address{
		Address: db.Email,
	}

	roles, err := role.ParseMany(db.Roles)
	if err != nil {
		return userbus.User{}, fmt.Errorf("parse: %w", err)
	}

	nme, err := name.Parse(db.Name)
	if err != nil {
		return userbus.User{}, fmt.Errorf("parse name: %w", err)
	}

	bus := userbus.User{
		ID:           db.ID,
		UserName:     nme,
		UserEmail:    addr,
		Roles:        roles,
		PasswordHash: db.PasswordHash,
		CreatedAt:    db.CreatedAt.In(time.Local),
	}

	return bus, nil
}
