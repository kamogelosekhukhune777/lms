// Package coursedb contains product related CRUD functionality.
package coursedb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

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

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
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

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, cor coursebus.CourseSchema) error {

	//insert course
	query := `
	INSERT INTO courses (
		ID, instructorId, instructorName, date, title, category, level, primaryLanguage,
		subtitle, description, image, welcomeMessage, pricing, objectives, isPublished
	) VALUES (:ID, :instructorId, :instructorName, :date, :title, :category, :level, :primaryLanguage,
		:subtitle, :description, :image, :welcomeMessage, :pricing, :objectives, :isPublished)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, query, toDBCourseSchema(cor)); err != nil {
		return fmt.Errorf("named exec context: %w", err)
	}

	//insert Lecture and link Lecture To Course
	if len(cor.Curriculum) > 0 {
		for _, lecture := range cor.Curriculum {

			query = `
			INSERT INTO lectures (lecture_id, title, video_url, public_id, free_preview)
			VALUES (:lecture_id, :title, :video_url, :public_id, :free_preview)`

			data := toDBLecture(lecture)
			data.ID = uuid.New()

			if err := sqldb.NamedExecContext(ctx, s.log, s.db, query, data); err != nil {
				return fmt.Errorf("named exec context: %w", err)
			}

			//link Lecture To Course
			query = `
			INSERT INTO curriculum (curriculum_id, course_id, lecture_id)
			VALUES ($1, $2, $3)`

			curriculumData := struct {
				CurriculumID uuid.UUID `db:"curriculum_id"`
				CourseID     uuid.UUID `db:"course_id"`
				LectureID    uuid.UUID `db:"lecture_id"`
			}{
				CurriculumID: uuid.New(),
				CourseID:     cor.ID,
				LectureID:    data.ID,
			}

			if err := sqldb.NamedExecContext(ctx, s.log, s.db, query, curriculumData); err != nil {
				return fmt.Errorf("namedexeccontext: %w", err)
			}
		}
	}

	//insert the students
	if len(cor.Students) > 0 {
		for _, student := range cor.Students {

			query = `
	        INSERT INTO course_students (course_student_id, course_id, student_id, student_name, student_email, paid_amount)
			VALUES ($1, $2, $3, $4, $5, $6)`

			data := toDBStudent(student)
			data.CourseID = cor.ID
			data.ID = uuid.New()

			if err := sqldb.NamedExecContext(ctx, s.log, s.db, query, data); err != nil {
				return fmt.Errorf("named exec context: %w", err)
			}
		}
	}

	return nil
}

func (s *Store) Update(ctx context.Context, cor coursebus.CourseSchema) error {
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

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCourseSchema(cor)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	//========================================================================================
	// Delete existing lectures and add updated lectures
	if len(cor.Curriculum) > 0 {
		data := cor.ID

		const q = `DELETE FROM curriculum WHERE course_id = $1`
		if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil { //????
			return fmt.Errorf("namedexeccontext: %w", err)
		}

		const qs = `DELETE FROM lectures WHERE lecture_id IN (
			SELECT lecture_id FROM curriculum WHERE course_id = $1
		)`
		if err := sqldb.NamedExecContext(ctx, s.log, s.db, qs, data); err != nil { //????
			return fmt.Errorf("namedexeccontext: %w", err)
		}

		if len(cor.Curriculum) > 0 {
			for _, lecture := range cor.Curriculum {

				const query = `
				INSERT INTO 
					lectures (lecture_id, title, video_url, public_id, free_preview)
				VALUES ($1, $2, $3, $4, $5)`

				data := toDBLecture(lecture)
				data.ID = uuid.New()

				if err := sqldb.NamedExecContext(ctx, s.log, s.db, query, data); err != nil {
					return fmt.Errorf("named exec context: %w", err)
				}

				//link Lecture To Course
				const q = `
				INSERT INTO
					curriculum (curriculum_id, course_id, lecture_id)
				VALUES ($1, $2, $3)`

				curriculumData := struct {
					CurriculumID uuid.UUID `db:"curriculum_id"`
					CourseID     uuid.UUID `db:"course_id"`
					LectureID    uuid.UUID `db:"lecture_id"`
				}{
					CurriculumID: uuid.New(),
					CourseID:     cor.ID,
					LectureID:    data.ID,
				}

				if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, curriculumData); err != nil {
					return fmt.Errorf("namedexeccontext: %w", err)
				}
			}
		}
	}

	if len(cor.Students) > 0 {
		// Delete existing students and add updated students
		const qd = `DELETE FROM course_students WHERE course_id = $1`
		data := cor.ID
		if err := sqldb.NamedExecContext(ctx, s.log, s.db, qd, data); err != nil {
			return fmt.Errorf("namedexeccontext: %w", err)
		}

		for _, student := range cor.Students {

			const query = `
				INSERT INTO 
					course_students (course_student_id, course_id, student_id, student_name, student_email, paid_amount)
				VALUES ($1, $2, $3, $4, $5, $6)`

			data := toDBStudent(student)
			data.CourseID = cor.ID
			data.ID = uuid.New()

			if err := sqldb.NamedExecContext(ctx, s.log, s.db, query, data); err != nil {
				return fmt.Errorf("named exec context: %w", err)
			}
		}
	}

	return nil
}

// filers
func (s *Store) QueryAll(ctx context.Context) ([]coursebus.CourseSchema, error) {
	const q = `
		SELECT 
			course_id, instructor_id, instructor_name, title, category, level, primary_language,
		    subtitle, description, image, welcome_message, pricing, objectives, is_published
		FROM 
			courses`

	var dbCourses []courseSchema

	if err := sqldb.QuerySlice(ctx, s.log, s.db, q, &dbCourses); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return []coursebus.CourseSchema{}, fmt.Errorf("db: %w", coursebus.ErrNotFound)
		}
		return []coursebus.CourseSchema{}, fmt.Errorf("db: %w", err)
	}

	return toBusCoursesSchema(dbCourses)
}

// Function: Get all student courses
func (s *Store) GetAllStudentCourses(ctx context.Context, filter coursebus.QueryFilter, orderBy order.By, page page.Page) ([]coursebus.CourseSchema, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
		SELECT 
			course_id, instructor_id, instructor_name, title, category, level, primary_language,
		    subtitle, description, image, welcome_message, pricing, objectives, is_published
		FROM 
			courses`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbCourses []courseSchema
	if err := sqldb.QuerySlice(ctx, s.log, s.db, buf.String(), &dbCourses); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return []coursebus.CourseSchema{}, fmt.Errorf("db: %w", coursebus.ErrNotFound)
		}
		return []coursebus.CourseSchema{}, fmt.Errorf("db: %w", err)
	}

	return toBusCoursesSchema(dbCourses) //?????
}

func (s *Store) QueryByID(ctx context.Context, courseID uuid.UUID) (coursebus.CourseSchema, error) {
	data := struct {
		ID string `db:"course_id"`
	}{
		ID: courseID.String(),
	}

	const q = `
		SELECT 
			course_id, instructor_id, instructor_name, title, category, level, primary_language,
		    subtitle, description, image, welcome_message, pricing, objectives, is_published
		FROM 
			courses
		WHERE 
			course_id = $1`

	var dbCor courseSchema
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbCor); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return coursebus.CourseSchema{}, fmt.Errorf("db: %w", coursebus.ErrNotFound)
		}
		return coursebus.CourseSchema{}, fmt.Errorf("db: %w", err)
	}

	//======================================================================================================
	// Fetch students

	const qs = `
	    SELECT 
			student_id, student_name, student_email, paid_amount
		FROM 
			course_students
		WHERE 
			course_id = $1`

	var dbStd []student
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, qs, data, &dbStd); err != nil {
		return coursebus.CourseSchema{}, fmt.Errorf("db: (fetchStudents) %w", err)
	}

	//======================================================================================================
	// Fetch lectures

	const ql = `
		SELECT 
			l.lecture_id, l.title, l.video_url, l.public_id, l.free_preview
		FROM 
			lectures l
		JOIN 
			curriculum c ON l.lecture_id = c.lecture_id
		WHERE 
			c.course_id = $1`

	var dbLec []lecture
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, ql, data, &dbLec); err != nil {
		return coursebus.CourseSchema{}, fmt.Errorf("db: (fetchStudents) %w", err)
	}

	// Map database records to business logic structures
	course, err := toBusCourseSchema(dbCor)
	if err != nil {
		return coursebus.CourseSchema{}, err
	}

	curriculum, err := toBusLectures(dbLec)
	if err != nil {
		return coursebus.CourseSchema{}, err
	}

	courseStudents, err := toBusStudents(dbStd)
	if err != nil {
		return coursebus.CourseSchema{}, err
	}

	// Attach curriculum and students to the course
	course.Curriculum = curriculum
	course.Students = courseStudents

	return course, nil
}

// //???????

func (s *Store) GetCurrentCourseProgress(ctx context.Context, userID, courseID uuid.UUID) (coursebus.CourseProgress, error) {
	data := struct {
		UserID   uuid.UUID `db:"user_id"`
		CourseID uuid.UUID `db:"course_id"`
	}{
		UserID:   userID,
		CourseID: courseID,
	}

	//========================================================================================================

	// Fetch the course progress
	const qs = `
		SELECT 
			completed, completion_date
		FROM 
			course_progress
		WHERE 
			user_id = $1 AND course_id = $2`

	type completionData struct {
		Completed      bool      `db:"completed"`
		CompletionDate time.Time `db:"completion_date"`
	}

	var dest completionData

	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, qs, data, &dest); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return coursebus.CourseProgress{}, fmt.Errorf("db: %w", coursebus.ErrNotFound)
		}
		return coursebus.CourseProgress{}, fmt.Errorf("db: %w", err)
	}

	//========================================================================================================
	// Fetch lecture progress
	const ql = `
		SELECT 
			lp.lecture_id, l.title, lp.viewed, lp.date_viewed
		FROM 
			lectures_progress lp
		JOIN 
			lectures l ON lp.lecture_id = l.lecture_id
		WHERE 
			lp.user_id = $1 AND lp.course_id = $2`

	var dbLecPro []lectureProgress

	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, ql, data, &dbLecPro); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return coursebus.CourseProgress{}, fmt.Errorf("db: %w", coursebus.ErrNotFound)
		}
		return coursebus.CourseProgress{}, fmt.Errorf("db: %w", err)
	}

	busLec := toBusLectureProgress(dbLecPro)

	//========================================================================================================

	dbCorPro := coursebus.CourseProgress{
		UserID:           userID,
		CourseID:         courseID,
		Completed:        dest.Completed,
		CompletionDate:   dest.CompletionDate,
		LecturesProgress: busLec,
	}

	return dbCorPro, nil
}

// complete
func (s *Store) MarkLectureAsViewed(ctx context.Context, userID, courseID, lectureID uuid.UUID) error {
	data := struct {
		DateViwed time.Time `db:"date_viewed"`
		UserID    uuid.UUID `db:"user_id"`
		CourseID  uuid.UUID `db:"course_id"`
		LectureID uuid.UUID `db:"lecture_id"`
	}{
		UserID:    userID,
		CourseID:  courseID,
		LectureID: lectureID,
		DateViwed: time.Now(),
	}

	const q = `
		UPDATE 
			CourseProgressLectures
		SET 
			viewed = TRUE, date_viewed = $1
		WHERE course_progress_id = (
			SELECT id FROM CourseProgress WHERE user_id = $2 AND course_id = $3
		) AND lecture_id = $4`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// complete
func (s *Store) ResetCourseProgress(ctx context.Context, userID, courseID uuid.UUID) error {
	data := struct {
		UserID   uuid.UUID `db:"user_id"`
		CourseID uuid.UUID `db:"course_id"`
	}{
		UserID:   userID,
		CourseID: courseID,
	}

	// Reset course completion
	const q = `
		UPDATE 
			course_progress
		SET 
			completed = FALSE, completion_date = NULL
		WHERE 
			user_id = $1 AND course_id = $2`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	// Reset lectures progress
	const ql = `
		UPDATE 
			lectures_progress
		SET 
			viewed = FALSE, date_viewed = NULL
		WHERE 
			course_progress_id = (
			SELECT course_progress_id FROM course_progress WHERE user_id = $1 AND course_id = $2)
		`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, ql, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
