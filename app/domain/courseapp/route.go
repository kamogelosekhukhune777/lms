package courseapp

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kamogelosekhukhune777/lms/app/sdk/mid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/sqldb"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log       *logger.Logger
	CourseBus *coursebus.Business
	UserBus   *userbus.Business
	DB        *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	cor := mid.GetCourseByID(cfg.CourseBus)
	usr := mid.GetUserByID(cfg.UserBus)
	transaction := mid.BeginCommitRollback(cfg.Log, sqldb.NewBeginner(cfg.DB))

	api := newApp(cfg.CourseBus, cfg.UserBus)

	//instructor
	app.HandlerFunc(http.MethodPost, version, "/add", api.create, transaction)
	app.HandlerFunc(http.MethodGet, version, "/get/details/{course_id}", api.queryByID, cor, transaction)
	app.HandlerFunc(http.MethodGet, version, "/get", api.queryAll, transaction)
	app.HandlerFunc(http.MethodPut, version, "/update/{course_id}", api.update, cor, transaction)

	//student routes
	//-course
	app.HandlerFunc(http.MethodGet, version, "/get", api.getAllStudentViewCourses, transaction)
	app.HandlerFunc(http.MethodGet, version, "get/details/{course_id}", api.getStudentViewCourseDetails, cor, transaction)                  //"/get/details/:id"
	app.HandlerFunc(http.MethodGet, version, "/purchase-info/{course_id}/{student_id}", api.checkCoursePurchaseInfo, usr, cor, transaction) //"/purchase-info/:id/:studentId""

	//-student-courses
	app.HandlerFunc(http.MethodGet, version, "/get/{user_id}", api.getCoursesByStudentId, usr, transaction) //----"/get/{student_id}"

	//-course progress
	app.HandlerFunc(http.MethodGet, version, "/get/{user_id}/{course_id}", api.getCurrentCourseProgress, usr, cor, transaction) //"get/:userId/:courseId"
	app.HandlerFunc(http.MethodPost, version, "/mark-lecture-viewed", api.markLectureAsViewed, transaction)
	app.HandlerFunc(http.MethodPost, version, "/reset-progress", api.resetCurrentCourseProgress, transaction)

}
