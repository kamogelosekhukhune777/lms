package studentbus

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/domain/userbus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/sqldb"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	CheckCoursePurchaseInfo(ctx context.Context, userID, courseID uuid.UUID) (bool, error)
	GetStudentCoursesByID(ctx context.Context, userID uuid.UUID) ([]Course, error)
	//GetAllStudentCourses(ctx context.Context) ([]StudentCourses, error)
	GetStudentByID(ctx context.Context, userID uuid.UUID) (StudentCourses, error)
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
}

// Business manages the set of APIs for product access.
type Business struct {
	log       *logger.Logger
	userBus   *userbus.Business
	courseBus *coursebus.Business
	storer    Storer
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

func (s *Business) CheckCoursePurchaseInfo(ctx context.Context, userID, courseID uuid.UUID) (bool, error) {

	return false, nil
}

// complete
func (s *Business) GetStudentCoursesByID(ctx context.Context, userID uuid.UUID) (StudentCourses, error) {
	return StudentCourses{}, nil
}

func (s *Business) GetAllStudentCourses(ctx context.Context) ([]StudentCourses, error) {

	return []StudentCourses{}, nil
}

func (s *Business) GetStudentByID(ctx context.Context, userID uuid.UUID) (StudentCourses, error) {
	return StudentCourses{}, nil
}
