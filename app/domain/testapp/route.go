package testapp

import (
	"net/http"

	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
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
	app.HandlerFunc(http.MethodGet, version, "/testerror", api.testError)
	app.HandlerFunc(http.MethodGet, version, "/testpanic", api.testPanic)
}
