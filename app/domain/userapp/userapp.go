// Package userapp maintains the app layer api for the user domain.
package userapp

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kamogelosekhukhune777/lms/app/sdk/auth"
	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

type app struct {
	userBus *userbus.Business
	auth    *auth.Auth
}

func newApp(userBus *userbus.Business, auth *auth.Auth) *app {
	return &app{
		userBus: userBus,
		auth:    auth,
	}
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	var app NewUser
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nc, err := toBusNewUser(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return errs.New(errs.Aborted, userbus.ErrUniqueEmail)
		}
		return errs.Newf(errs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	claims := auth.Claims{
		Roles: []string{"USER"},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    a.auth.Issuer(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(100 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token, err := a.auth.GenerateToken(claims)
	if err != nil {
		return errs.Newf(errs.Internal, "create: failed to generate token: %s", err)
	}

	return toAppUserWithToken(&usr, token)

}

func (a *app) logIn(ctx context.Context, r *http.Request) web.Encoder {
	var app logInUser
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	email, err := mail.ParseAddress(app.Email)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	// Find user by email
	usr, err := a.userBus.Authenticate(ctx, *email, app.Password)
	if err != nil {
		if errors.Is(err, userbus.ErrAuthenticationFailure) {
			return errs.New(errs.Unauthenticated, errors.New("invalid email or password"))
		}
		return errs.Newf(errs.Internal, "logIn: failed to authenticate user: %s", err)
	}

	claims := auth.Claims{
		Roles: []string{"USER"},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    a.auth.Issuer(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(100 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token, err := a.auth.GenerateToken(claims)
	if err != nil {
		return errs.Newf(errs.Internal, "logIn: failed to generate token: %s", err)
	}

	return toAppUserWithToken(&usr, token)
}

func (a *app) checkAuth(ctx context.Context, r *http.Request) web.Encoder {
	return nil
}
