package userapp

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/app/sdk/auth"
	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

type app struct {
	userBus *userbus.Business
	auth    *auth.Auth
}

func newApp(userBus *userbus.Business, ath *auth.Auth) *app {
	return &app{
		userBus: userBus,
		auth:    ath,
	}
}

// creates a user
func (a *app) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

	err = token(usr, w, a.auth)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusOK)
}

func (a *app) login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app logInUser
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	lu, err := tologInUser(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Authenticate(ctx, lu.UserEmail, lu.Password)
	if err != nil {
		return errs.Newf(errs.Internal, "")
	}

	err = token(usr, w, a.auth)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, lu, http.StatusOK)
	//return web.Respond(ctx,w,usr,http.StatusOk)
}

func (a *app) checkauth(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := a.auth.Authenticate(ctx, r.Header.Get("authorization"))
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	if claims.Subject == "" {
		return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims")
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return errs.New(errs.Unauthenticated, errors.New("ID is not in its proper form"))
	}

	user, err := a.userBus.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, userbus.ErrNotFound):
			return errs.New(errs.Unauthenticated, err)
		default:
			return errs.Newf(errs.Unauthenticated, "querybyid: userID[%s]: %s", id, err)
		}
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Authenticated user!",
		"data": map[string]interface{}{
			"user": user,
		},
	}

	return web.Respond(ctx, w, response, http.StatusOK)
}
