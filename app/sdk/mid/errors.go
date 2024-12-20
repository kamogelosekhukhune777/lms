package mid

import (
	"context"
	"errors"
	"net/http"
	"path"

	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

// Errors handles errors coming out of the call chain.
func Errors(log *logger.Logger) web.MidFunc {
	m := func(next web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			err := next(ctx, w, r)
			if err == nil {
				return nil
			}

			var appErr *errs.Error
			if !errors.As(err, &appErr) {
				appErr = errs.Newf(errs.Internal, "Internal Server Error")

			}
			log.Error(ctx, "handled error during request",
				"err", err,
				"source_err_file", path.Base(appErr.FileName),
				"source_err_func", path.Base(appErr.FuncName))

			if appErr.Code == errs.InternalOnlyLog {
				appErr = errs.Newf(errs.Internal, "Internal Server Error")
			}

			// Send the error to the transport package so the error can be
			// used as the response.

			return appErr

		}

		return h
	}

	return m
}
