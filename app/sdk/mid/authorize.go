package mid

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

// ErrInvalidID represents a condition where the id is not a uuid.
var ErrInvalidID = errors.New("ID is not in its proper form")

// AuthorizeCourse executes the specified role and extracts the specified
// user from the DB if a user id is specified in the call. Depending on the rule
// specified, the userid from the claims may be compared with the specified
// user id.
func AuthorizeCourse(courseBus *coursebus.Business) web.MidFunc {
	m := func(next web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			id := web.Param(r, "course_id")

			var courseID uuid.UUID

			if id != "" {
				var err error
				courseID, err = uuid.Parse(id)
				if err != nil {
					return errs.New(errs.Unauthenticated, ErrInvalidID)
				}

				cor, err := courseBus.QueryByID(ctx, courseID)
				if err != nil {
					switch {
					case errors.Is(err, coursebus.ErrNotFound):
						return errs.New(errs.Unauthenticated, err)
					default:
						return errs.Newf(errs.Unauthenticated, "querybyid: courseID[%s]: %s", courseID, err)
					}
				}

				ctx = setCourse(ctx, cor)
			}

			return next(ctx, w, r)
		}

		return h
	}

	return m
}
