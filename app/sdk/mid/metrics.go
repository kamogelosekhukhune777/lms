package mid

import (
	"context"
	"net/http"

	"github.com/kamogelosekhukhune777/vendly/app/sdk/metrics"
	"github.com/kamogelosekhukhune777/vendly/foundation/web"
)

// Metrics updates program counters.
func Metrics() web.MidFunc {
	m := func(next web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = metrics.Set(ctx)

			err := next(ctx, w, r)

			n := metrics.AddRequests(ctx)

			if n%1000 == 0 {
				metrics.AddGoroutines(ctx)
			}

			if err != nil {
				metrics.AddErrors(ctx)
			}

			return err
		}

		return h
	}

	return m
}
