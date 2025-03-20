package orderdb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/orderbus"
	"github.com/kamogelosekhukhune777/lms/business/types/money"
)

type order struct {
	ID            uuid.UUID `db:"order_id"`
	UserID        uuid.UUID `db:"user_id"`
	OrderStatus   string    `db:"order_status"`
	PaymentMethod string    `db:"payment_method"`
	PaymentStatus string    `db:"payment_status"`
	OrderDate     time.Time `db:"order_date"`
	PaymentID     string    `db:"payment_id"`
	PayerID       string    `db:"payer_id"`
	InstructorID  uuid.UUID `db:"instructor_id"`
	CourseID      uuid.UUID `db:"course_id"`
	CoursePricing float64   `db:"course_pricing"`
}

func toDBOrder(ord orderbus.Order) order {
	return order{
		ID:            ord.ID,
		UserID:        ord.UserID,
		OrderStatus:   ord.OrderStatus,
		PaymentMethod: ord.PaymentMethod,
		PaymentStatus: ord.PaymentStatus,
		OrderDate:     ord.OrderDate.UTC(),
		PaymentID:     ord.PaymentID,
		PayerID:       ord.PayerID,
		InstructorID:  ord.InstructorID,
		CourseID:      ord.CourseID,
		CoursePricing: ord.CoursePricing.Value(),
	}
}

func toBusOrder(db order) (orderbus.Order, error) {
	price, err := money.Parse(db.CoursePricing)
	if err != nil {
		return orderbus.Order{}, fmt.Errorf("parse cost: %w", err)
	}

	bus := orderbus.Order{
		ID:            db.ID,
		UserID:        db.UserID,
		OrderStatus:   db.OrderStatus,
		PaymentMethod: db.PaymentMethod,
		PaymentStatus: db.PaymentStatus,
		OrderDate:     db.OrderDate.In(time.Local),
		PaymentID:     db.PaymentID,
		PayerID:       db.PayerID,
		InstructorID:  db.InstructorID,
		CourseID:      db.CourseID,
		CoursePricing: price,
	}

	return bus, nil
}
