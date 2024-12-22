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

	/*

		authen := mid.Authenticate(cfg.AuthClient)
		ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)
		ruleAuthorizeUser := mid.AuthorizeUser(cfg.AuthClient, cfg.UserBus, auth.RuleAdminOrSubject)
	*/

	api := newApp(cfg.UserBus, cfg.Auth)

	app.HandlerFunc(http.MethodPost, version, "/auth/registerUser", api.create)
	app.HandlerFunc(http.MethodPost, version, "/auth/loginUser", api.login)
	app.HandlerFunc(http.MethodPost, version, "/auth/check-auth", api.checkauth)

}
