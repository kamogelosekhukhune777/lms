package mid

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/kamogelosekhukhune777/vendly/app/sdk/errs"
	"github.com/kamogelosekhukhune777/vendly/foundation/web"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func Panics() web.MidFunc {
	m := func(next web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (resp error) {

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					resp = errs.Newf(errs.InternalOnlyLog, "PANIC [%v] TRACE[%s]", rec, string(trace))

					//metrics.AddPanics(ctx)
				}
			}()

			return next(ctx, w, r)
		}

		return h
	}

	return m
}
