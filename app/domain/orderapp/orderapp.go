package orderapp

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/app/sdk/paypal"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/domain/orderbus"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

type app struct {
	orderBus  *orderbus.Business
	courseBus *coursebus.Business
	userBus   *userbus.Business
	paypal    *paypal.PayPalClient
}

func newApp(courseBus *coursebus.Business, userBus *userbus.Business, paypal *paypal.PayPalClient) *app {
	return &app{
		courseBus: courseBus,
		userBus:   userBus,
		paypal:    paypal,
	}
}

// CreateOrder handles creating a PayPal order
func (a *app) createOrder(ctx context.Context, r *http.Request) web.Encoder {
	var app NewOrder
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	order, err := toBusNewOrder(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	paypal, err := a.paypal.CreateOrder(ctx, order.CoursePricing.String(), "USD")
	if err != nil {
		return errs.New(errs.Internal, errors.New("failed to create payapal order"))
	}

	order.PaymentID = paypal.ID
	order.OrderStatus = "pending"
	order.PaymentStatus = "pending"

	ord, err := a.orderBus.SaveOrder(ctx, order)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return toAppOrder(ord)
}

// CapturePayment handles capturing a PayPal payment
func (a *app) capturePayment(ctx context.Context, r *http.Request) web.Encoder {
	var request requestData

	if err := web.Decode(r, &request); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	orderID, err := uuid.Parse(request.OrderID)
	if err != nil {
		return errs.New(errs.Internal, errors.New("OrderID is not in its proper form"))
	}

	ord, err := a.orderBus.GetOrderByID(ctx, orderID)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	_, err = a.paypal.CaptureOrder(request.PaymentID)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	ord.OrderStatus = "confirmed"
	ord.PaymentStatus = "paid"

	if err := a.orderBus.UpdateOrder(ctx, ord); err != nil {
		return errs.New(errs.Internal, err)
	}

	return toAppOrder(ord)
}
