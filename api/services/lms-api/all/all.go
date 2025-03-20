// Package all binds all the routes into the specified app.
package all

import (
	"github.com/kamogelosekhukhune777/lms/app/domain/courseapp"
	"github.com/kamogelosekhukhune777/lms/app/domain/testapp"
	"github.com/kamogelosekhukhune777/lms/app/domain/userapp"
	"github.com/kamogelosekhukhune777/lms/app/sdk/mux"
	"github.com/kamogelosekhukhune777/lms/app/sdk/paypal"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg mux.Config) {

	testapp.Routes(app)

	userapp.Routes(app, userapp.Config{
		Log:     cfg.Log,
		UserBus: cfg.BusConfig.UserBus,
		Auth:    cfg.Auth,
	})

	courseapp.Routes(app, courseapp.Config{
		Log:       cfg.Log,
		CourseBus: cfg.BusConfig.CourseBus,
		UserBus:   cfg.BusConfig.UserBus,
		DB:        cfg.DB,
	})

	_, _ = paypal.NewPayPalClient("", "", "")
}
