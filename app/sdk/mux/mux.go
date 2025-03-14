// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kamogelosekhukhune777/lms/app/sdk/auth"
	"github.com/kamogelosekhukhune777/lms/app/sdk/mid"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

type BusConfig struct {
	UserBus *userbus.Business
}

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build     string
	Log       *logger.Logger
	DB        *sqlx.DB
	Auth      *auth.Auth
	BusConfig BusConfig
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg Config)
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config, routeAdder RouteAdder) http.Handler {
	logger := func(ctx context.Context, msg string, args ...any) {
		cfg.Log.Info(ctx, msg, args...)
	}

	app := web.NewApp(
		logger,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Panics(),
	)

	routeAdder.Add(app, cfg)

	return app
}
