package mid

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kamogelosekhukhune777/vendly/foundation/logger"
	"github.com/kamogelosekhukhune777/vendly/foundation/web"
)

// Logger writes information about the request to the logs.
func Logger(log *logger.Logger) web.MidFunc {
	m := func(next web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			now := time.Now()

			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Info(ctx, "request started", "method", r.Method, "path", path, "remoteaddr", r.RemoteAddr)

			err := next(ctx, w, r)

			log.Info(ctx, "request completed", "method", r.Method, "path", path, "remoteaddr", r.RemoteAddr,
				"since", time.Since(now).String())

			return err
		}

		return h
	}

	return m
}
