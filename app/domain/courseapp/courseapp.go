// Package productapp maintains the app layer api for the product domain.
package courseapp

import (
	"context"
	"net/http"

	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/app/sdk/mid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

type app struct {
	courseBus *coursebus.Business
}

func newApp(courseBus *coursebus.Business) *app {
	return &app{
		courseBus: courseBus,
	}
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	var app NewCourse
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	np, err := toBusNewCourse(ctx, app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	prd, err := a.courseBus.Create(ctx, np)
	if err != nil {
		return errs.Newf(errs.Internal, "create: prd[%+v]: %s", prd, err)
	}

	return toAppCourse(prd)
}

func (a *app) update(ctx context.Context, r *http.Request) web.Encoder {
	var app UpdateCourse
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	up, err := toBusUpdateCourse(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	prd, err := mid.GetCourse(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "product missing in context: %s", err)
	}

	updPrd, err := a.courseBus.Update(ctx, prd, up)
	if err != nil {
		return errs.Newf(errs.Internal, "update: productID[%s] up[%+v]: %s", prd.ID, app, err)
	}

	return toAppCourse(updPrd)
}

func (a *app) queryAll(ctx context.Context, r *http.Request) web.Encoder {
	cors, err := a.courseBus.QueryAll(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	return toAppCourses(cors)
}

func (a *app) queryByID(ctx context.Context, r *http.Request) web.Encoder {
	prd, err := mid.GetCourse(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "product missing in context: %s", err)
	}

	return toAppCourse(prd)
}

//==================================================================================================================

// ==================================================================================================================

func (a *app) getCoursesByStudentId(ctx context.Context, r *http.Request) web.Encoder {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "user missing in context: %s", err)
	}

	cors, err := a.courseBus.GetCoursesByStudentID(ctx, usr.ID)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	return toAppCourses(cors)
}
