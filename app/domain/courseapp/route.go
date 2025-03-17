package courseapp

import (
	"net/http"

	"github.com/kamogelosekhukhune777/lms/app/sdk/mid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log       *logger.Logger
	CourseBus *coursebus.Business
	UserBus   *userbus.Business
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	cor := mid.GetCourseByID(cfg.CourseBus)
	usr := mid.GetUserByID(cfg.UserBus)
	api := newApp(cfg.CourseBus)

	//instructor
	app.HandlerFunc(http.MethodPost, version, "/add", api.create)
	app.HandlerFunc(http.MethodGet, version, "/get/details/{course_id}", api.queryByID, cor)
	app.HandlerFunc(http.MethodGet, version, "/get", api.queryAll)
	app.HandlerFunc(http.MethodPut, version, "/update/{course_id}", api.update, cor)

	//student
	//course
	app.HandlerFunc(http.MethodGet, version, "/get", api.getAllStudentViewCourses)
	app.HandlerFunc(http.MethodGet, version, "get/details/{course_id}", api.getStudentViewCourseDetails, cor)                  //"/get/details/:id"
	app.HandlerFunc(http.MethodGet, version, "/purchase-info/{course_id}/{student_id}", api.checkCoursePurchaseInfo, usr, cor) //"/purchase-info/:id/:studentId""

	//student-courses
	app.HandlerFunc(http.MethodGet, version, "/get/{user_id}", api.getCoursesByStudentId, usr) //----"/get/{student_id}"

	//course progress

}
