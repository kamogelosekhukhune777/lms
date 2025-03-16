package courseapp

import (
	"net/http"

	"github.com/kamogelosekhukhune777/lms/app/sdk/mid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log       *logger.Logger
	CourseBus *coursebus.Business
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	cor := mid.GetCourseByID(cfg.CourseBus)
	api := newApp(cfg.CourseBus)

	//instructor
	app.HandlerFunc(http.MethodPost, version, "/add", api.create)
	app.HandlerFunc(http.MethodPost, version, "/get/details/{course_id}", api.queryByID, cor)
	app.HandlerFunc(http.MethodPost, version, "/get", api.queryAll)
	app.HandlerFunc(http.MethodPost, version, "/update/{course_id}", api.update, cor)

	//student
	//course
}
