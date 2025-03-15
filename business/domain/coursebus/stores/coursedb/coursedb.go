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

//=============================================================================================================================
