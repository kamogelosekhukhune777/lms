package courseapp

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kamogelosekhukhune777/lms/app/sdk/mid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/sqldb"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

type Config struct {
	DB        *sqlx.DB
	Log       *logger.Logger
	CourseBus *coursebus.Business
}

func Routes(app web.App, cfg Config) {
	const version = "v1"

	api := newApp(cfg.CourseBus)
	transaction := mid.BeginCommitRollback(cfg.Log, sqldb.NewBeginner(cfg.DB))
	authorizeCourse := mid.AuthorizeCourse(cfg.CourseBus)

	// Courses Management
	app.HandlerFunc(http.MethodPost, version, "/add", api.create, transaction)
	app.HandlerFunc(http.MethodGet, version, "/get", api.getAllCourses)
	app.HandlerFunc(http.MethodPut, version, "/get/details/:id", api.getCourseDetails, authorizeCourse)
	app.HandlerFunc(http.MethodPut, version, "/update/:id", api.update, authorizeCourse, transaction)

	// Course Progress
	app.HandlerFunc(http.MethodGet, version, "/get/:userId/:courseId", api.getCurrentCourseProgress)
	app.HandlerFunc(http.MethodPost, version, "/mark-lecture-viewed", api.markCurrentLectureAsViewed)
	app.HandlerFunc(http.MethodPost, version, "/reset-progress", api.resetCurrentCourseProgress)

}
