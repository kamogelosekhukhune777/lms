// Package all binds all the routes into the specified app.
package all

import (
	"github.com/kamogelosekhukhune777/lms/app/domain/testapp"
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

	testapp.Routes(app, testapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
	})
}
