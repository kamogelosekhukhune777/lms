package userdb

import (
	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/sqldb/dbarray"
	"github.com/kamogelosekhukhune777/lms/business/types/role"
)

type user struct {
	ID           uuid.UUID      `db:"user_id"`
	Name         string         `db:"user_name"`
	Email        string         `db:"user_email"`
	PasswordHash []byte         `db:"password_hash"`
	Roles        dbarray.String `db:"roles"`
}

func toDBUser(bus userbus.User) user {
	return user{
		ID:           bus.ID,
		Name:         bus.UserName.String(),
		Email:        bus.UserEmail.Address,
		Roles:        role.ParseToString(bus.Roles),
		PasswordHash: bus.PasswordHash,
	}
}
