package studentapp

import (
	"context"
	"net/http"

	"github.com/kamogelosekhukhune777/lms/app/sdk/mid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/domain/studentbus"
)

type app struct {
	studentBus *studentbus.Business
	courseBus  *coursebus.Business
}

func newApp(courseBus *coursebus.Business, studentBus *studentbus.Business) *app {
	return &app{
		studentBus: studentBus,
		courseBus:  courseBus,
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

	studentBus, err := a.studentBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	app := app{
		studentBus: studentBus,
		courseBus:  courseBus,
	}

	return &app, nil
}

func (a *app) getAllStudentViewCourses(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (a *app) getStudentViewCourseDetails(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (a *app) checkCoursePurchaseInfo(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}
