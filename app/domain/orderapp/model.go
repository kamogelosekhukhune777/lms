package orderapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/business/domain/orderbus"
	"github.com/kamogelosekhukhune777/lms/business/types/money"
)

type Order struct {
	//ID             string
	UserID         string    `json:"userId"`
	UserName       string    `json:"userName"`
	UserEmail      string    `json:"userEmail"`
	OrderStatus    string    `json:"orderStatus"`
	PaymentMethod  string    `json:"paymentMethod"`
	PaymentStatus  string    `josn:"paymentStatus"`
	OrderDate      time.Time `json:"orderDate"`
	PaymentId      string    `json:"paymentId"`
	PayerId        string    `json:"payerId"`
	InstructorId   string    `json:"instructorId"`
	InstructorName string    `json:"instructorName"`
	CourseImage    string    `json:"courseImage"`
	CourseTitle    string    `json:"courseTitle"`
	CourseId       string    `json:"courseId"`
	CoursePricing  string    `json:"coursePricing"`
}

// Encode implements the encoder interface.
func (app Order) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppOrder(ord orderbus.Order) Order {
	return Order{
		OrderDate: ord.OrderDate,
	}
}

//============================================================================

type NewOrder struct {
	UserID         string    `json:"userId" validate:"required"`
	UserName       string    `json:"userName" validate:"required"`
	UserEmail      string    `json:"userEmail" validate:"required"`
	OrderStatus    string    `json:"orderStatus" validate:"required"`
	PaymentMethod  string    `json:"paymentMethod" validate:"required"`
	PaymentStatus  string    `josn:"paymentStatus" validate:"required"`
	OrderDate      time.Time `json:"orderDate" validate:"required"`
	PaymentID      string    `json:"paymentId" validate:"required"`
	PayerID        string    `json:"payerId" validate:"required"`
	InstructorID   string    `json:"instructorId" validate:"required"`
	InstructorName string    `json:"instructorName" validate:"required"`
	CourseImage    string    `json:"courseImage" validate:"required"`
	CourseTitle    string    `json:"courseTitle" validate:"required"`
	CourseID       string    `json:"courseId" validate:"required"`
	CoursePricing  float64   `json:"coursePricing" validate:"required,gte=0"`
}

// Decode implements the decoder interface.
func (app *NewOrder) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewOrder) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewOrder(app NewOrder) (orderbus.NewOrder, error) {
	userID, err := uuid.Parse(app.UserID)
	if err != nil {
		return orderbus.NewOrder{}, errs.New(errs.Internal, errors.New("ID is not in its proper form"))
	}

	instructorID, err := uuid.Parse(app.UserID)
	if err != nil {
		return orderbus.NewOrder{}, errs.New(errs.Internal, errors.New("ID is not in its proper form"))
	}

	courseID, err := uuid.Parse(app.UserID)
	if err != nil {
		return orderbus.NewOrder{}, errs.New(errs.Internal, errors.New("ID is not in its proper form"))
	}

	price, err := money.Parse(app.CoursePricing)
	if err != nil {
		return orderbus.NewOrder{}, fmt.Errorf("parse cost: %w", err)
	}

	bus := orderbus.NewOrder{
		UserID:         userID,
		UserName:       app.UserName,
		UserEmail:      app.UserEmail,
		OrderStatus:    app.OrderStatus,
		PaymentMethod:  app.PaymentMethod,
		PaymentStatus:  app.PaymentStatus,
		OrderDate:      app.OrderDate,
		PaymentID:      app.PaymentID,
		PayerID:        app.PaymentID,
		InstructorID:   instructorID,
		InstructorName: app.InstructorName,
		CourseImage:    app.CourseImage,
		CourseTitle:    app.CourseTitle,
		CourseID:       courseID,
		CoursePricing:  price,
	}

	return bus, nil
}

//=======================================================================

type requestData struct {
	PaymentID string `json:"paymentId"`
	PayerID   string `json:"payerId"`
	OrderID   string `json:"orderId"`
}

// Decode implements the decoder interface.
func (app *requestData) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app requestData) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}
