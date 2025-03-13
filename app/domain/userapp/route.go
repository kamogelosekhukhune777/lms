package userapp

import (
	"net/http"

	"github.com/kamogelosekhukhune777/lms/app/sdk/auth"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log     *logger.Logger
	UserBus *userbus.Business
	Auth    *auth.Auth
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newApp(cfg.UserBus, cfg.Auth)

	app.HandlerFunc(http.MethodGet, version, "/check-auth", api.checkAuth)
	app.HandlerFunc(http.MethodPost, version, "/register", api.create)
	app.HandlerFunc(http.MethodPut, version, "/login", api.logIn)
}
