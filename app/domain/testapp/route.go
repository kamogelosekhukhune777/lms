package testapp

import (
	"net/http"

	"github.com/kamogelosekhukhune777/vendly/foundation/logger"
	"github.com/kamogelosekhukhune777/vendly/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newApp(cfg.Build, cfg.Log)

	app.HandlerFunc(http.MethodGet, version, "/test", api.test)
}
