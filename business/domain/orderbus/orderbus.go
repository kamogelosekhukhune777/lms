// Package corsebus provides business access to product domain.
package orderbus

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/sqldb"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, no Order) error
	Update(ctx context.Context, ord Order) error
	QueryByID(ctx context.Context, orderID uuid.UUID) (Order, error)
}

// Business manages the set of APIs for product access.
type Business struct {
	log       *logger.Logger
	userBus   *userbus.Business
	courseBus *coursebus.Business
	storer    Storer
}

// NewBusiness constructs a product business API for use.
func NewBusiness(log *logger.Logger, userBus *userbus.Business, courseBus *coursebus.Business, storer Storer) *Business {
	b := Business{
		log:       log,
		userBus:   userBus,
		courseBus: courseBus,
		storer:    storer,
	}

	return &b
}

// NewWithTx constructs a new business value that will use the
// specified transaction in any store related calls.
func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (*Business, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	userBus, err := b.userBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	courseBus, err := b.courseBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	bus := Business{
		log:       b.log,
		userBus:   userBus,
		courseBus: courseBus,
		storer:    storer,
	}

	return &bus, nil
}

// SaveOrder saves the order in the database (mocked)
func (b *Business) SaveOrder(ctx context.Context, no NewOrder) (Order, error) {

	usr, err := b.userBus.QueryByID(ctx, no.UserID)
	if err != nil {
		return Order{}, fmt.Errorf("user.querybyid: %s: %w", no.UserID, err)
	}

	cor, err := b.courseBus.QueryByID(ctx, no.CourseID)
	if err != nil {
		return Order{}, fmt.Errorf("course.querybyid: %s: %w", no.CourseID, err)
	}

	order := Order{
		ID:             uuid.New(),
		UserID:         usr.ID,
		UserName:       no.UserName,
		UserEmail:      no.UserEmail,
		OrderStatus:    no.OrderStatus,
		PaymentMethod:  no.PaymentMethod,
		PaymentStatus:  no.PaymentStatus,
		OrderDate:      no.OrderDate,
		PaymentID:      no.PaymentID,
		PayerID:        no.PaymentID,
		InstructorID:   no.InstructorID,
		InstructorName: no.InstructorName,
		CourseImage:    no.CourseImage,
		CourseTitle:    no.CourseTitle,
		CourseID:       cor.ID,
		CoursePricing:  no.CoursePricing,
	}

	if err := b.storer.Create(ctx, order); err != nil {
		return Order{}, fmt.Errorf("create: %w", err)
	}

	return order, nil
}

// GetOrderByID retrieves an order from the database (mocked)
func (b *Business) GetOrderByID(ctx context.Context, orderID uuid.UUID) (Order, error) {
	prd, err := b.storer.QueryByID(ctx, orderID)
	if err != nil {
		return Order{}, fmt.Errorf("query: orderID[%s]: %w", orderID, err)
	}

	return prd, nil
}

// UpdateOrder updates an order in the database (mocked)
func (b *Business) UpdateOrder(ctx context.Context, ord Order) error {
	// Mock database update
	if err := b.storer.Update(ctx, ord); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil

}

// EnrollStudentInCourse updates student records (mocked)
func (b *Business) EnrollStudentInCourse(userID uuid.UUID, order *Order) error {
	// Mock update student enrollment
	return nil
}
