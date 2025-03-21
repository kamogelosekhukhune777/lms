// Package all binds all the routes into the specified app.
package all

import (
	"github.com/kamogelosekhukhune777/lms/app/domain/courseapp"
	"github.com/kamogelosekhukhune777/lms/app/domain/mediapp"
	"github.com/kamogelosekhukhune777/lms/app/domain/orderapp"
	"github.com/kamogelosekhukhune777/lms/app/domain/testapp"
	"github.com/kamogelosekhukhune777/lms/app/domain/userapp"
	"github.com/kamogelosekhukhune777/lms/app/sdk/mux"
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

	orderapp.Routes(app, orderapp.Config{
		Log:       cfg.Log,
		CourseBus: cfg.BusConfig.CourseBus,
		UserBus:   cfg.BusConfig.UserBus,
		Paypal:    cfg.Paypal,
	})

	mediapp.Routes(app, mediapp.Config{
		Log:              cfg.Log,
		CloudinaryClient: cfg.CloudinaryClient,
	})
}
