package studentapp

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kamogelosekhukhune777/lms/app/sdk/mid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/domain/studentbus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/sqldb"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

type Config struct {
	DB         *sqlx.DB
	Log        *logger.Logger
	CourseBus  *coursebus.Business
	StudentBus *studentbus.Business
}

func Routes(app web.App, cfg Config) {
	const version = "v1"

	api := newApp(cfg.CourseBus, cfg.StudentBus)
	transaction := mid.BeginCommitRollback(cfg.Log, sqldb.NewBeginner(cfg.DB))

	//Student View Courses
	app.HandlerFunc(http.MethodGet, version, "/get", api.getAllStudentViewCourses, transaction)
	app.HandlerFunc(http.MethodGet, version, "/get/details/:id", api.getStudentViewCourseDetails, transaction)
	app.HandlerFunc(http.MethodGet, version, "/purchase-info/:id/:studentId", api.checkCoursePurchaseInfo, transaction)
}
