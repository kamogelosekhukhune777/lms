package testapp

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
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

func (a *app) testError(ctx context.Context, r *http.Request) web.Encoder {
	if n := rand.Intn(100); n%2 == 0 {
		return errs.Newf(errs.FailedPrecondition, "this message is trused")
	}

	app := Test{
		Status: "OK",
	}

	return app
}

func (a *app) testPanic(ctx context.Context, r *http.Request) web.Encoder {

	if n := rand.Intn(100); n%2 == 0 {
		panic("WE ARE PANICKING!!!")
	}

	app := Test{
		Status: "OK",
	}

	return app
}
