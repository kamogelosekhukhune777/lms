// Package productapp maintains the app layer api for the product domain.
package courseapp

import (
	"context"
	"net/http"

	"fmt"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/app/sdk/mid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/order"
	"github.com/kamogelosekhukhune777/lms/business/sdk/page"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

type app struct {
	courseBus *coursebus.Business
	userBus   *userbus.Business
}

func newApp(courseBus *coursebus.Business, userBus *userbus.Business) *app {
	return &app{
		courseBus: courseBus,
		userBus:   userBus,
	}
}

// newWithTx constructs a new Handlers value with the domain apis
// using a store transaction that was created via middleware.
func (a *app) newWithTx(ctx context.Context) (*app, error) {
	tx, err := mid.GetTran(ctx)
	if err != nil {
		return nil, err
	}

	userBus, err := a.userBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	courseBus, err := a.courseBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	app := app{
		userBus:   userBus,
		courseBus: courseBus,
	}

	return &app, nil
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	var app NewCourse
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
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

	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
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
	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	cors, err := a.courseBus.QueryAll(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	return toAppCourses(cors)
}

func (a *app) queryByID(ctx context.Context, r *http.Request) web.Encoder {
	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	prd, err := mid.GetCourse(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "product missing in context: %s", err)
	}

	return toAppCourse(prd)
}

//==================================================================================================================

func (a *app) getAllStudentViewCourses(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseQueryParams(r)

	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, coursebus.DefaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	prds, err := a.courseBus.GetAllStudentViewCourses(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	return toAppCourses(prds)
}

func (a *app) getStudentViewCourseDetails(ctx context.Context, r *http.Request) web.Encoder {
	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	cor, err := mid.GetCourse(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "course missing in context: %s", err)
	}

	lecs, err := a.courseBus.GetLectures(ctx, cor.ID)
	if err != nil {
		return errs.Newf(errs.Internal, "%s", err)
	}

	cor.Curriculum = lecs

	return toAppCourse(cor)
}

func (a *app) checkCoursePurchaseInfo(ctx context.Context, r *http.Request) web.Encoder {
	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	usr, err := mid.GetUser(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "user missing in context: %s", err)
	}

	cor, err := mid.GetCourse(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "course missing in context: %s", err)
	}

	sta, err := a.courseBus.CheckCoursePurchaseInfo(ctx, cor.ID, usr.ID)
	if err != nil {
		return errs.Newf(errs.Internal, "purchase info : %s", err)
	}

	return BoolResult(sta)
}

// ==================================================================================================================

func (a *app) getCoursesByStudentId(ctx context.Context, r *http.Request) web.Encoder {
	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

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

//=====================================================================================================================

func (a *app) getCurrentCourseProgress(ctx context.Context, r *http.Request) web.Encoder {
	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	//unecessay
	cor, err := mid.GetCourse(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "course missing in context: %s", err)
	}

	usr, err := mid.GetUser(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "user missing in context: %s", err)
	}

	fmt.Println(cor, usr)

	return nil
}

func (a *app) markLectureAsViewed(ctx context.Context, r *http.Request) web.Encoder {
	values := r.URL.Query()

	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	userID, err := uuid.Parse(values.Get("user_Id"))
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	courseID, err := uuid.Parse(values.Get("course_Id"))
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	lectureID, err := uuid.Parse(values.Get("lecture_Id"))
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	corp, err := a.courseBus.MarkLecture(ctx, userID, courseID, lectureID)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return toAppCourseProgress(corp)
}

func (a *app) resetCurrentCourseProgress(ctx context.Context, r *http.Request) web.Encoder {
	values := r.URL.Query()

	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	userID, err := uuid.Parse(values.Get("user_Id"))
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	courseID, err := uuid.Parse(values.Get("course_Id"))
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	corp, err := a.courseBus.ResetCourseProgress(ctx, userID, courseID)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return toAppCourseProgress(corp)
}
