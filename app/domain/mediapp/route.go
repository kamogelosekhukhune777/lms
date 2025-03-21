package mediapp

import (
	"net/http"

	"github.com/kamogelosekhukhune777/lms/app/sdk/cloudinary"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log              *logger.Logger
	CloudinaryClient *cloudinary.CloudinaryService
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newApp(cfg.CloudinaryClient)

	app.HandlerFunc(http.MethodPost, version, "/uplaod", api.uploadFile)
	app.HandlerFunc(http.MethodDelete, version, "/delete/{id}", api.deleteFile)
	app.HandlerFunc(http.MethodPost, version, "/bulk-upload", api.bulkUpload)
}
