package mid

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
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

			var statusCode = errs.OK
			if err != nil {
				statusCode = errs.Internal

				var v *errs.Error
				if errors.As(err, &v) {
					statusCode = v.Code
				}
			}

			log.Info(ctx, "request completed", "method", r.Method, "path", path, "remoteaddr", r.RemoteAddr,
				"statuscode", statusCode, "since", time.Since(now).String())

			return err
		}

		return h
	}

	return m
}
