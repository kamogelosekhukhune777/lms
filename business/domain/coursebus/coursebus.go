// Package coursebus provides business access to course domain.
package coursebus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/order"
	"github.com/kamogelosekhukhune777/lms/business/sdk/page"
	"github.com/kamogelosekhukhune777/lms/business/sdk/sqldb"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound    = errors.New("course not found")
	ErrInvalidCost = errors.New("cost not valid")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, cor Course) error
	Update(ctx context.Context, cor Course) error
	QueryByID(ctx context.Context, courseID uuid.UUID) (Course, error)
	QueryAll(ctx context.Context) ([]Course, error)
	GetCoursesByStudentID(ctx context.Context, studentId uuid.UUID) ([]Course, error)
	CheckCoursePurchaseInfo(ctx context.Context, courseID uuid.UUID, studentID uuid.UUID) (bool, error)
	GetLectures(ctx context.Context, courseID uuid.UUID) ([]Lecture, error)
	GetCoureStudents(ctx context.Context, courseID uuid.UUID) ([]Student, error)
	QueryAllStudentViewCourses(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Course, error)
}

// Business manages the set of APIs for product access.
type Business struct {
	log     *logger.Logger
	userBus *userbus.Business
	storer  Storer
}

// NewBusiness constructs a product business API for use.
func NewBusiness(log *logger.Logger, userBus *userbus.Business, storer Storer) *Business {
	b := Business{
		log:     log,
		userBus: userBus,
		storer:  storer,
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

	bus := Business{
		log:     b.log,
		userBus: userBus,
		storer:  storer,
	}

	return &bus, nil
}

// Create adds a new Course to the system.
func (b *Business) Create(ctx context.Context, np NewCourse) (Course, error) {

	/*
		image saving
	*/

	usr, err := b.userBus.QueryByID(ctx, np.InstructorID)
	if err != nil {
		return Course{}, fmt.Errorf("course.querybyid: %s: %w", np.InstructorID, err)
	}

	now := time.Now()

	cor := Course{
		ID:              uuid.New(),
		InstructorID:    usr.ID,
		Title:           np.Title,
		Category:        np.Category,
		Level:           np.Level,
		PrimaryLanguage: np.PrimaryLanguage,
		Subtitle:        np.Subtitle,
		Description:     np.Description,
		Image:           np.Image,
		WelcomeMessage:  np.WelcomeMessage,
		Pricing:         np.Pricing,
		Objectives:      np.Objectives,
		IsPublished:     true,
		CreatedAt:       now,
	}

	if err := b.storer.Create(ctx, cor); err != nil {
		return Course{}, fmt.Errorf("create: %w", err)
	}

	return cor, nil
}

// Update modifies information about a Course.
func (b *Business) Update(ctx context.Context, cor Course, upc UpdateCourse) (Course, error) {
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

	if upc.Image != nil {
		cor.Image = *upc.Image
	}

	if upc.WelcomeMessage != nil {
		cor.WelcomeMessage = *upc.WelcomeMessage
	}

	if upc.Pricing != nil {
		cor.Pricing = *upc.Pricing
	}

	if upc.Objectives != nil {
		cor.Objectives = *upc.Objectives
	}

	/*if upc.IsPublished != nil {
		cor.IsPublished = *upc.Ispublished
	}*/

	if err := b.storer.Update(ctx, cor); err != nil {
		return Course{}, fmt.Errorf("update: %w", err)
	}

	return cor, nil
}

// QueryByID finds the course by the specified ID.
func (b *Business) QueryByID(ctx context.Context, courseID uuid.UUID) (Course, error) {

	cor, err := b.storer.QueryByID(ctx, courseID)
	if err != nil {
		return Course{}, fmt.Errorf("query: courseID[%s]: %w", courseID, err)
	}

	return cor, nil
}

func (b *Business) QueryAll(ctx context.Context) ([]Course, error) {
	cors, err := b.storer.QueryAll(ctx)
	if err != nil {
		return []Course{}, fmt.Errorf("query: %w", err)
	}

	return cors, nil
}

//==================================================================================================================

func (b *Business) GetAllStudentViewCourses(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Course, error) {

	prds, err := b.storer.QueryAllStudentViewCourses(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil

}

func (b *Business) CheckCoursePurchaseInfo(ctx context.Context, courseID uuid.UUID, studentID uuid.UUID) (bool, error) {
	sta, err := b.storer.CheckCoursePurchaseInfo(ctx, courseID, studentID)
	if err != nil {
		return false, fmt.Errorf("query: %w", err)
	}

	return sta, nil
}

func (b *Business) GetLectures(ctx context.Context, courseID uuid.UUID) ([]Lecture, error) {
	lecs, err := b.storer.GetLectures(ctx, courseID)
	if err != nil {
		return []Lecture{}, nil
	}

	return lecs, nil
}

func (b *Business) GetCoureStudents(ctx context.Context, courseID uuid.UUID) ([]Student, error) {
	stu, err := b.storer.GetCoureStudents(ctx, courseID)
	if err != nil {
		return []Student{}, nil
	}

	return stu, nil
}

//==================================================================================================================

func (b *Business) GetCoursesByStudentID(ctx context.Context, studentId uuid.UUID) ([]Course, error) {
	cors, err := b.storer.GetCoursesByStudentID(ctx, studentId)
	if err != nil {
		return []Course{}, fmt.Errorf("query: %w", err)
	}

	return cors, nil
}
