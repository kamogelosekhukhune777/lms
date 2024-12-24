// Package coursedb contains product related CRUD functionality.
package coursedb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
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

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, cor coursebus.Course) error {
	const q = `
	INSERT INTO products
		(product_id, user_id, name, cost, quantity, date_created, date_updated)
	VALUES
		(:product_id, :user_id, :name, :cost, :quantity, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCourse(cor)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Update(ctx context.Context, cor coursebus.Course) error {
	const q = `
	UPDATE
		courses
	SET
		"instructor_name" = :instructor_name,
		"date" = :date,
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
		"is_published" = :is_published
	WHERE
		course_id = :course_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCourse(cor)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) QueryAll(ctx context.Context) ([]coursebus.Course, error) {
	const q = `
	SELECT 
		course_id,
		instructor_id,
		instructor_name,
		date,
		title,
		category,
		level,
		primary_language,
		subtitle,
		description,
		image,
		welcome_message,
		pricing,
		objectives,
		is_published
	FROM 
		courses
	WHERE 
		is_published = TRUE;`

	var dbCourses []course

	if err := sqldb.QuerySlice(ctx, s.log, s.db, q, &dbCourses); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return []coursebus.Course{}, fmt.Errorf("db: %w", coursebus.ErrNotFound)
		}
		return []coursebus.Course{}, fmt.Errorf("db: %w", err)
	}

	return toBusCourses(dbCourses)
}

func (s *Store) QueryByID(ctx context.Context, courseID uuid.UUID) (coursebus.Course, error) {
	data := struct {
		ID string `db:"course_id"`
	}{
		ID: courseID.String(),
	}

	const q = `
	SELECT
	    course_id, instructor_id, instructor_name, date, title, category, level, primary_language, subtitle, description, image, welcome_message, pricing, objectives, is_published
	FROM
		courses
	WHERE
		course_id_id = :course_id`

	var dbCor course
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbCor); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return coursebus.Course{}, fmt.Errorf("db: %w", coursebus.ErrNotFound)
		}
		return coursebus.Course{}, fmt.Errorf("db: %w", err)
	}

	return toBusCourse(dbCor)
}
