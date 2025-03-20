package orderapp

import (
	"net/http"

	"github.com/kamogelosekhukhune777/lms/app/sdk/paypal"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log       *logger.Logger
	courseBus *coursebus.Business
	userBus   *userbus.Business
	paypal    *paypal.PayPalClient
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newApp(cfg.courseBus, cfg.userBus, cfg.paypal)

	app.HandlerFunc(http.MethodPost, version, "/create", api.createOrder)
	app.HandlerFunc(http.MethodPost, version, "/capture", api.capturePayment) //capturePaymentAndFinalizeOrder
}
