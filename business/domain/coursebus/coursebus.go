// Package coursebus provides business access to product domain.
package coursebus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound = errors.New("course not found")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, cor Course) error
	Update(ctx context.Context, cor Course) error
	QueryAll(ctx context.Context) ([]Course, error)
	QueryByID(ctx context.Context, id uuid.UUID) (Course, error)
}

// Business manages the set of APIs for product access.
type Business struct {
	log     *logger.Logger
	userBus *userbus.Business
	//delegate *delegate.Delegate
	storer Storer
}

// NewBusiness constructs a product business API for use.
func NewBusiness(log *logger.Logger, userBus *userbus.Business, storer Storer) *Business {
	b := Business{
		log:     log,
		userBus: userBus,
		//delegate: delegate,
		storer: storer,
	}

	//b.registerDelegateFunctions()

	return &b
}

func (b *Business) NewCourse(ctx context.Context, nc NewCourse) (Course, error) {

	now := time.Now()

	//image saving
	//
	//

	cor := Course{
		ID:              uuid.New(),
		InstructorId:    nc.InstructorId,   //
		InstructorName:  nc.InstructorName, //
		Date:            now,
		Title:           nc.Title,
		Category:        nc.Category,
		Level:           nc.Level,
		PrimaryLanguage: nc.PrimaryLanguage,
		Subtitle:        nc.Subtitle,
		Description:     nc.Description,
		Pricing:         nc.Pricing,
		WelcomeMessage:  nc.WelcomeMessage,
		Image:           nc.Image,
	}

	if err := b.storer.Create(ctx, cor); err != nil {
		return Course{}, fmt.Errorf("create: %w", err)
	}

	return cor, nil
}

func (b *Business) GetAllCourses() ([]Course, error) {
	ctx := context.Background()

	courses, err := b.storer.QueryAll(ctx)
	if err != nil {
		return []Course{}, err
	}

	return courses, nil
}

func (b *Business) GetCourseDetailsByID(ctx context.Context, id uuid.UUID) (Course, error) {
	cor, err := b.storer.QueryByID(ctx, id)
	if err != nil {
		return Course{}, fmt.Errorf("query: productID[%s]: %w", id, err)
	}

	return cor, nil
}

func (b *Business) Update(ctx context.Context, cor Course, upc UpdateCousre) (Course, error) {

	if upc.Title != nil {
		cor.Title = *upc.Title
	}

	if upc.Category != nil {
		cor.Category = *upc.Category
	}

	if upc.Level != nil {
		cor.Level = *upc.Level
	}

	if upc.PrimaryLanguage != nil {
		cor.PrimaryLanguage = *upc.PrimaryLanguage
	}

	if upc.Subtitle != nil {
		cor.Subtitle = *upc.Subtitle
	}

	if upc.Description != nil {
		cor.Description = *upc.Description
	}

	if upc.Pricing != nil {
		cor.Pricing = *upc.Pricing
	}

	if upc.WelcomeMessage != nil {
		cor.WelcomeMessage = *upc.WelcomeMessage
	}

	if upc.Image != nil {
		cor.Image = *upc.Image
	}

	err := b.storer.Update(ctx, cor)
	if err != nil {
		return Course{}, fmt.Errorf("update: %w", err)
	}

	return cor, nil
}
