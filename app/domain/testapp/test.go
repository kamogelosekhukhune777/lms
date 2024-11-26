package testapp

import (
	"context"
	"net/http"

	"github.com/kamogelosekhukhune777/multi-vendor-ecom/foundation/logger"
	"github.com/kamogelosekhukhune777/multi-vendor-ecom/foundation/web"
)

type app struct {
	build string
	log   *logger.Logger
}

func newApp(build string, log *logger.Logger) *app {
	return &app{
		build: build,
		log:   log,
	}
}

func (a *app) test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
