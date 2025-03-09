package testapp

import (
	"context"
	"net/http"

	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

type app struct{}

func newApp() *app {
	return &app{}
}

func (a *app) test(ctx context.Context, r *http.Request) web.Encoder {
	app := Test{
		Status: "OK",
	}

	return app
}
