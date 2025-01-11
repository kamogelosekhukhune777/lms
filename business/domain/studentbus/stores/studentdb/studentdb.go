package studentdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/domain/studentbus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/sqldb"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
)

// Store manages the set of APIs for user database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (studentbus.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

// Function: Check course purchase info
func (s *Store) CheckCoursePurchaseInfo(ctx context.Context, userID, courseID uuid.UUID) (bool, error) {
	data := struct {
		UserID   uuid.UUID `db:"user_id"`
		CourseID uuid.UUID `db:"course_id"`
	}{
		UserID:   userID,
		CourseID: courseID,
	}

	var dbcourses []course

	const q = `
		SELECT 
			user_id, course_id, title, instructor_id, instructor_name, date_of_purchase, course_image
		FROM 
			student_courses
		WHERE 
			user_id = $1
	`

	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbcourses); err != nil {
		return false, err
	}

	for _, v := range dbcourses {
		if v.CourseID == courseID {
			return true, nil
		}
	}

	return false, nil
}

// complete
func (s *Store) GetStudentCoursesByID(ctx context.Context, userID uuid.UUID) ([]studentbus.Course, error) {
	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	const q = `
			SELECT 
				c.course_id, c.title, c.instructor_id, c.instructor_name, c.date_of_purchase, c.course_image
			FROM 
				student_courses c
			WHERE 
				c.user_id = $1`

	var dbCors []course

	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbCors); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return []studentbus.Course{}, fmt.Errorf("db: %w", coursebus.ErrNotFound)
		}
		return []studentbus.Course{}, fmt.Errorf("db: %w", err)
	}

	return toBusStudentCourses(dbCors)
}

// =======================================================================================================================

// Function: Get student by ID
func (s *Store) GetStudentByID(ctx context.Context, userID uuid.UUID) (studentbus.StudentCourses, error) {
	data := struct {
		userID uuid.UUID `db:"user_id"`
	}{
		userID: userID,
	}

	var dbcourses []course

	const q = `
		SELECT 
			user_id, course_id, title, instructor_id, instructor_name, date_of_purchase, course_image
		FROM 
			student_courses
		WHERE 
			user_id = $1
	`

	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbcourses); err != nil {
		return studentbus.StudentCourses{}, err
	}

	courses, err := toBusStudentCourses(dbcourses)
	if err != nil {
		return studentbus.StudentCourses{}, nil
	}

	return studentbus.StudentCourses{
		UserID:  userID,
		Courses: courses,
	}, nil
}
