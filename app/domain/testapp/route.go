package testapp

import (
	"net/http"

	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App) {
	const version = "v1"

	api := newApp()

	app.HandlerFunc(http.MethodGet, version, "/test", api.test)
}
