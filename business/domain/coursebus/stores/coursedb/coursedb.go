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
	UPDATE Courses
	SET
		title = :title,
		category = :category,
		level = :level,
		primary_language = :primary_language,
		subtitle = :subtitle,
		description = :description,
		image = :image,
		welcome_message = :welcome_message,
		pricing = :pricing, 
		objectives = :objectives,
		is_published = :is_published,
		created_at = :created_at
	WHERE
		course_id = :course_id`

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

	var hasPurchased bool
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &hasPurchased); err != nil {
		return false, fmt.Errorf("namedquerystruct: %w", err)
	}

	return hasPurchased, nil
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
		CourseID string `db:"course_id"`
	}{
		CourseID: courseID.String(),
	}

	const query = `
	SELECT 
		s.* 
	FROM 
		students s
	JOIN enrollments e ON s.id = e.student_id
	WHERE 
		e.course_id = :course_id
	ORDER BY s.id`

	var students []student
	err := sqldb.NamedQuerySlice(ctx, s.log, s.db, query, data, &students)
	if err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
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

//====================================================================================================================

func (s *Store) GetCourseProgress(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (coursebus.CourseProgress, error) {
	data := struct {
		UserID   string `db:"user_id"`
		CourseID string `db:"course_id"`
	}{
		UserID:   userID.String(),
		CourseID: courseID.String(),
	}

	const q = `
	SELECT * 
	FROM 
		CourseProgress 
	WHERE 
		user_id = :user_id 
	AND 
		course_id = :course_id`

	var corp courseProgress
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &corp); err != nil {
		return coursebus.CourseProgress{}, fmt.Errorf("db: %w", err)
	}

	return toBusCourseProgress(corp), nil
}

func (s *Store) MarkLectureAsViewed(ctx context.Context, userID, courseID, lectureID uuid.UUID) error {
	data := struct {
		ID         string `db:"lecture_progress_id"`
		ProgressID string `db:"progress_id"`
		CourseID   string `db:"course_id"`
		UserID     string `db:"user_id"`
		LectureID  string `db:"lecture_id"`
	}{
		ID:         uuid.New().String(),
		ProgressID: uuid.New().String(),
		CourseID:   courseID.String(),
		UserID:     userID.String(),
		LectureID:  lectureID.String(),
	}

	// Mark lecture as viewed
	const ql = `
		INSERT INTO LectureProgress 
			(lecture_progress_id, user_id, lecture_id, viewed, date_viewed)
		VALUES 
			(:lecture_progress_id, :user_id, :lecture_id, TRUE, NOW())
		ON CONFLICT (user_id, lecture_id) DO UPDATE 
			SET viewed = TRUE, date_viewed = NOW()`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, ql, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	// Check if all lectures in the course are viewed by the user
	const qm = `
	SELECT NOT EXISTS 
		(SELECT 1 FROM Lectures l
		LEFT JOIN LectureProgress lp 
		ON l.lecture_id = lp.lecture_id AND lp.user_id = :user_id
		WHERE l.course_id = :course_id 
		AND (lp.viewed IS NULL OR lp.viewed = FALSE)) AS all_completed`

	var allCompleted bool
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, qm, data, &allCompleted); err != nil {
		return fmt.Errorf("namedquerystruct: %w", err)
	}

	// If all lectures are viewed, mark the course as completed
	if allCompleted {
		const qo = `
		INSERT INTO CourseProgress 
			(progress_id, user_id, course_id, completed, completion_date)
		VALUES 
			(:progress_id, :user_id, :course_id, TRUE, NOW())
		ON CONFLICT (user_id, course_id) DO UPDATE 
			SET completed = TRUE, completion_date = NOW()`

		if err := sqldb.NamedExecContext(ctx, s.log, s.db, qo, data); err != nil {
			return fmt.Errorf("namedexeccontext: %w", err)
		}
	}

	return nil
}

func (s *Store) ResetCourseProgress(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) error {
	data := struct {
		CourseID string `db:"course_id"`
		UserID   string `db:"user_id"`
	}{
		CourseID: courseID.String(),
		UserID:   userID.String(),
	}

	const ql = `
		DELETE FROM LectureProgress 
		WHERE 
			user_id = :user_id
		AND lecture_id IN 
			(SELECT lecture_id FROM Lectures WHERE course_id = :course_id)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, ql, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	// Delete course progress
	const qp = `
		DELETE FROM CourseProgress 
		WHERE 
			user_id = :user_id 
		AND 
			course_id = :course_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, qp, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
