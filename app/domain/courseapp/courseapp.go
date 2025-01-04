// Package courseapp provides business access to product domain.
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

// newWithTx constructs a new Handlers value with the domain apis
// using a store transaction that was created via middleware.
func (a *app) newWithTx(ctx context.Context) (*app, error) {
	tx, err := mid.GetTran(ctx)
	if err != nil {
		return nil, err
	}

	courseBus, err := a.courseBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	app := app{
		courseBus: courseBus,
	}

	return &app, nil
}

func (a *app) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewCourse
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	nc, err := toBusNewCourse(ctx, app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	cor, err := a.courseBus.Create(ctx, nc)
	if err != nil {
		return errs.Newf(errs.Internal, "create: cor[%+v]: %s", cor, err)
	}

	return web.Respond(ctx, w, toAppCourse(cor), http.StatusOK)
}

func (a *app) getAllCourses(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	cors, err := a.courseBus.GetAllCourses()
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	return web.Respond(ctx, w, toAppCourses(cors), http.StatusOK)
}

func (a *app) getCourseDetails(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	cor, err := mid.GetCourse(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "course missing in context: %s", err)
	}

	return web.Respond(ctx, w, cor, http.StatusOK)
}

func (a *app) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app UpdateCourse
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	up, err := toBusUpdateCourse(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	cor, err := mid.GetCourse(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "course missing in context: %s", err)
	}

	updCor, err := a.courseBus.Update(ctx, cor, up)
	if err != nil {
		return errs.Newf(errs.Internal, "update: courseID[%s] up[%+v]: %s", cor.ID, app, err)
	}

	return web.Respond(ctx, w, toAppCourse(updCor), http.StatusOK)
}
