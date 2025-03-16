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

func GetCourseByID(courseBus *coursebus.Business) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			id := web.Param(r, "course_id")

			if id != "" {
				var err error
				courseID, err := uuid.Parse(id)
				if err != nil {
					return errs.New(errs.Unauthenticated, ErrInvalidID)
				}

				prd, err := courseBus.QueryByID(ctx, courseID)
				if err != nil {
					switch {
					case errors.Is(err, coursebus.ErrNotFound):
						return errs.New(errs.Unauthenticated, err)
					default:
						return errs.Newf(errs.Internal, "querybyid: courseID[%s]: %s", courseID, err)
					}
				}

				ctx = setCourse(ctx, prd)
			}

			return next(ctx, r)
		}

		return h
	}

	return m
}
