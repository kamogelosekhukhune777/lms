package mid

import (
	"context"
	"net/http"

	"github.com/kamogelosekhukhune777/lms/app/sdk/metrics"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

// Metrics updates program counters.
func Metrics() web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			ctx = metrics.Set(ctx)

			resp := next(ctx, r)

			n := metrics.AddRequests(ctx)

			if n%1000 == 0 {
				metrics.AddGoroutines(ctx)
			}

			if isError(resp) != nil {
				metrics.AddErrors(ctx)
			}

			return resp
		}

		return h
	}

	return m
}
