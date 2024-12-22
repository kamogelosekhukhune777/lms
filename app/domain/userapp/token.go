package userapp

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kamogelosekhukhune777/lms/app/sdk/auth"
	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
)

func token(usr userbus.User, w http.ResponseWriter, a *auth.Auth) error {
	roles := []string{}
	for _, v := range usr.Roles {
		roles = append(roles, v.String())
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    a.Issuer(),
			Subject:   usr.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(100 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: roles,
	}

	tkn, err := a.GenerateToken(claims)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	w.Header().Set("Authorization", "Bearer "+tkn)

	return nil
}
