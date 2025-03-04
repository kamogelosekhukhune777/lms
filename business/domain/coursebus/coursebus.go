// Package coursebus provides business access to product domain.
package coursebus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/sqldb"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound = errors.New("course not found")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, cor CourseSchema) error
	Update(ctx context.Context, cor CourseSchema) error
	QueryAll(ctx context.Context) ([]CourseSchema, error)
	QueryByID(ctx context.Context, id uuid.UUID) (CourseSchema, error)
	GetCurrentCourseProgress(ctx context.Context, userID, courseID uuid.UUID) (CourseProgress, error)
	MarkLectureAsViewed(ctx context.Context, userID, courseID, lectureID uuid.UUID) error
	ResetCourseProgress(ctx context.Context, userID, courseID uuid.UUID) error
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

// NewWithTx constructs a new business value that will use the
// specified transaction in any store related calls.
func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (*Business, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	bus := Business{
		log:    b.log,
		storer: storer,
	}

	return &bus, nil
}

func (b *Business) Create(ctx context.Context, nc NewCourseSchema) (CourseSchema, error) {

	now := time.Now()

	//image saving and video saving
	//
	//

	cor := CourseSchema{
		ID:              uuid.New(),
		InstructorID:    nc.InstructorID,   //
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
		Curriculum:      nc.Curriculum,
		Objectives:      nc.Objectives,
	}

	if err := b.storer.Create(ctx, cor); err != nil {
		return CourseSchema{}, fmt.Errorf("create: %w", err)
	}

	return cor, nil
}

func (b *Business) GetAllCourses() ([]CourseSchema, error) {
	ctx := context.Background()

	courses, err := b.storer.QueryAll(ctx)
	if err != nil {
		return []CourseSchema{}, err
	}

	return courses, nil
}

func (b *Business) QueryByID(ctx context.Context, id uuid.UUID) (CourseSchema, error) {
	cor, err := b.storer.QueryByID(ctx, id)
	if err != nil {
		return CourseSchema{}, fmt.Errorf("query: productID[%s]: %w", id, err)
	}

	return cor, nil
}

func (b *Business) Update(ctx context.Context, cor CourseSchema, upc UpdateCourseSchema) (CourseSchema, error) {

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

	if upc.Students != nil {
		if !isNilOrEmpty(upc.Students) && !slicesEqualUnordered(cor.Students, upc.Students) {
			cor.Students = upc.Students
		}
	}

	if upc.Curriculum != nil {
		if !isNilOrEmpty(upc.Curriculum) && !slicesEqualUnordered(cor.Curriculum, upc.Curriculum) {
			cor.Curriculum = upc.Curriculum
		}
	}

	err := b.storer.Update(ctx, cor)
	if err != nil {
		return CourseSchema{}, fmt.Errorf("update: %w", err)
	}

	return cor, nil
}

func (b *Business) GetCurrentCourseProgress(ctx context.Context, userID, courseID uuid.UUID) (CourseProgress, error) {
	data, err := b.storer.GetCurrentCourseProgress(ctx, userID, courseID)
	if err != nil {
		return CourseProgress{}, nil
	}

	return data, nil
}

func (b *Business) ResetCourseProgress(ctx context.Context, userID, courseID uuid.UUID) error {
	err := b.storer.ResetCourseProgress(ctx, userID, courseID)
	if err != nil {
		return err
	}

	return nil
}

func (b *Business) MarkLectureAsViewed(ctx context.Context, userID, courseID, lectureID uuid.UUID) error {
	err := b.storer.MarkLectureAsViewed(ctx, userID, courseID, lectureID)
	if err != nil {
		return err
	}

	return nil
}
