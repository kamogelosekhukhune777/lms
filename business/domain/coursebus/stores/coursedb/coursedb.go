package coursedb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/order"
	"github.com/kamogelosekhukhune777/lms/business/sdk/page"
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

func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (coursebus.Storer, error) {
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

func (s *Store) Create(ctx context.Context, cor coursebus.Course) error {
	const q = `
	INSERT INTO Courses
		(course_id, instructor_id, title, category, level, primary_language, subtitle, description, image, welcome_message, pricing, objectives, is_published, created_at)
	VALUES
		(:course_id, :instructor_id, :title, :category, :level, :primary_language, :subtitle, :description, :image, :welcome_message, :pricing, :objectives, :is_published, :created_at)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCourse(cor)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Update(ctx context.Context, cor coursebus.Course) error {
	const q = `
	UPDATE
		Courses
	SET
		"title" = :title,
		"category" = :category,
		"level" = :level,
		"primary_language" = :primary_language,
		"subtitle" = :subtitle,
		"description" = :description,
		"image" = :image,
		"welcome_message" = :welcome_message,
		"pricing" = :pricing, 
		"objectives" = :objectives,
		"is_published" = :is_published,
		"created_at" = :created_at
	WHERE
		product_id = :product_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCourse(cor)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) QueryByID(ctx context.Context, courseID uuid.UUID) (coursebus.Course, error) {
	data := struct {
		ID string `db:"course_id"`
	}{
		ID: courseID.String(),
	}

	const q = `
	SELECT
	    course_id, instructor_id, title, category, level, primary_language, subtitle, description, image, welcome_message, pricing, objectives, is_published, created_at
	FROM
		Courses
	WHERE
		course_id = :course_id`

	var dbPrd course
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbPrd); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return coursebus.Course{}, fmt.Errorf("db: %w", coursebus.ErrNotFound)
		}
		return coursebus.Course{}, fmt.Errorf("db: %w", err)
	}

	return toBusCourse(dbPrd)

}

func (s *Store) QueryAll(ctx context.Context) ([]coursebus.Course, error) {
	const q = `
	SELECT
	    course_id, instructor_id, title, category, level, primary_language, subtitle, description, image, welcome_message, pricing, objectives, is_published, created_at
	FROM
		Courses`

	var dbPrds []course
	if err := sqldb.QuerySlice(ctx, s.log, s.db, q, &dbPrds); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusCourses(dbPrds)
}

//==============================================================================================================================

func (s *Store) QueryAllStudentViewCourses(ctx context.Context, filter coursebus.QueryFilter, orderBy order.By, page page.Page) ([]coursebus.Course, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
	    product_id, user_id, name, cost, quantity, date_created, date_updated
	FROM
		Courses`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbPrds []course
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbPrds); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusCourses(dbPrds)
}

func (s *Store) CheckCoursePurchaseInfo(ctx context.Context, courseID uuid.UUID, studentID uuid.UUID) (bool, error) {
	data := struct {
		courseID  string `db:"course_id"`
		StudentID string `db:"student_id"`
	}{
		courseID:  courseID.String(),
		StudentID: studentID.String(),
	}

	const q = `
	SELECT EXISTS (
		SELECT 1 
		FROM Enrollments 
		WHERE student_id = :student_id 
		AND course_id = :course_id
	) AS has_purchased`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return false, nil
	}

	return true, nil
}

func (s *Store) GetLectures(ctx context.Context, courseID uuid.UUID) ([]coursebus.Lecture, error) {
	data := struct {
		CourseID string `json:"course_id"`
	}{
		CourseID: courseID.String(),
	}

	query := `
	SELECT 
		id, title 
	FROM 
		lectures 
	WHERE 
		course_id = :course_id
	ORDER BY id`

	var lectures []lecture
	err := sqldb.NamedQuerySlice(ctx, s.log, s.db, query, data, &lectures)
	if err != nil {
		return nil, err
	}
	return toBusLectures(lectures)
}

func (s *Store) GetCoureStudents(ctx context.Context, courseID uuid.UUID) ([]coursebus.Student, error) {

	data := struct {
		CourseID string `json:"course_id"`
	}{
		CourseID: courseID.String(),
	}

	query := `
	SELECT 
		s.id, s.name 
	FROM 
		students s
	JOIN enrollments e ON s.id = e.student_id
	WHERE 
		e.course_id = :course_id
	ORDER BY s.id`

	var students []student
	err := sqldb.NamedQuerySlice(ctx, s.log, s.db, query, data, &students)
	if err != nil {
		return nil, err
	}

	return toBusStudents(students)
}

// =============================================================================================================================

func (s *Store) GetCoursesByStudentID(ctx context.Context, studentID uuid.UUID) ([]coursebus.Course, error) {

	data := struct {
		ID string `db:"user_id"`
	}{
		ID: studentID.String(),
	}

	const q = `
	SELECT 
		c.course_id,
		c.title,
		c.category,
		c.level,
		c.primary_language,
		c.subtitle,
		c.description,
		c.image,
		c.welcome_message,
		c.pricing,
		c.objectives,
		c.is_published,
		c.created_at,
		e.enrolled_at
	FROM Enrollments e
		JOIN Courses c ON e.course_id = c.course_id
	WHERE e.student_id = :user_id`

	var dbPrds []course
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbPrds); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusCourses(dbPrds)
}
