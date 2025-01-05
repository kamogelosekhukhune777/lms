// Package coursedb contains product related CRUD functionality.
package coursedb

import (
	"context"
	"errors"
	"fmt"
	"time"

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
func (s *Store) Create(ctx context.Context, cor coursebus.Course) error {

	//insert course
	query := `
	INSERT INTO courses (
		ID, instructorId, instructorName, date, title, category, level, primaryLanguage,
		subtitle, description, image, welcomeMessage, pricing, objectives, isPublished
	) VALUES (:ID, :instructorId, :instructorName, :date, :title, :category, :level, :primaryLanguage,
		:subtitle, :description, :image, :welcomeMessage, :pricing, :objectives, :isPublished)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, query, toDBCourse(cor)); err != nil {
		return fmt.Errorf("named exec context: %w", err)
	}

	//insert the students
	if len(cor.Students) > 0 {
		for _, student := range cor.Students {
			query = `
	           INSERT INTO Students (courseId, studentId, studentName, studentEmail, paidAmount)
	           VALUES (:courseId, :studentId, :studentName, :studentEmail, :paidAmount)`

			data := toDBStudent(student)

			if err := sqldb.NamedExecContext(ctx, s.log, s.db, query, data); err != nil {
				return fmt.Errorf("named exec context: %w", err)
			}
		}
	}

	//insert Lecture and link Lecture To Course
	if len(cor.Curriculum) > 0 {
		for _, lecture := range cor.Curriculum {
			query = `
			INSERT INTO lectures (lecture_id, title, videoUrl, public_id, freePreview)
			VALUES (:lecture_id, :title, :videoUrl, :public_id, :freePreview)`

			data := toDBLecture(lecture)

			if err := sqldb.NamedExecContext(ctx, s.log, s.db, query, data); err != nil {
				return fmt.Errorf("named exec context: %w", err)
			}

			//link Lecture To Course
			query = `
			INSERT INTO curriculum (course_id, lecture_id)
			VALUES (:course_id, :lecture_id)`

			curriculumData := toDBCurriculum(lecture.ID, cor.ID)

			if err := sqldb.NamedExecContext(ctx, s.log, s.db, query, curriculumData); err != nil {
				return fmt.Errorf("namedexeccontext: %w", err)
			}
		}
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

	//========================================================================================
	// Delete existing lectures and add updated lectures
	if len(cor.Curriculum) > 0 {
		const qd = `
			DELETE FROM 
				CourseLectures 
			WHERE 
				course_id = $1`

		if err := sqldb.ExecContext(ctx, s.log, s.db, qd); err != nil { //????
			return fmt.Errorf("namedexeccontext: %w", err)
		}

		for _, lecture := range cor.Curriculum {
			const q = `
				INSERT INTO Lectures (
					title, video_url, public_id, free_preview
				) VALUES ($1, $2, $3, $4)`

			if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBLecture(lecture)); err != nil {
				return fmt.Errorf("named exec context: %w", err)
			}

			const qs = `
				INSERT INTO lectures (course_id, lecture_id)
				SELECT $1, id FROM Lectures WHERE public_id = $2`

			curriculumData := toDBCurriculum(lecture.ID, cor.ID)

			if err := sqldb.NamedExecContext(ctx, s.log, s.db, qs, curriculumData); err != nil {
				return fmt.Errorf("namedexeccontext: %w", err)
			}
		}
	}

	if len(cor.Students) > 0 {
		// Delete existing students and add updated students
		const qd = `DELETE FROM CourseStudents WHERE course_id = $1`
		if err := sqldb.ExecContext(ctx, s.log, s.db, qd); err != nil { //????
			return fmt.Errorf("namedexeccontext: %w", err)
		}

		for _, student := range cor.Students {
			const query = `
	           INSERT INTO Students (courseId, studentId, studentName, studentEmail, paidAmount)
	           VALUES (:courseId, :studentId, :studentName, :studentEmail, :paidAmount)`

			if err := sqldb.NamedExecContext(ctx, s.log, s.db, query, toDBStudent(student)); err != nil {
				return fmt.Errorf("named exec context: %w", err)
			}
		}
	}

	return nil
}

func (s *Store) QueryAll(ctx context.Context) ([]coursebus.Course, error) {
	const q = `
		SELECT 
			id, instructor_id, instructor_name, date, title, category, level, primary_language, subtitle, description, image, welcome_message, pricing, objectives, is_published 
		FROM 
			courses`

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
            id, instructor_id, instructor_name, date, title, category, level, primary_language, 
            subtitle, description, image, welcome_message, pricing, objectives, is_published
        FROM 
            courses 
        WHERE 
            id = $1`

	var dbCor course
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbCor); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return coursebus.Course{}, fmt.Errorf("db: %w", coursebus.ErrNotFound)
		}
		return coursebus.Course{}, fmt.Errorf("db: %w", err)
	}

	//======================================================================================================
	// Fetch students

	const qs = `
	    SELECT
            student_id, student_name, student_email, paid_amount 
        FROM 
            students 
        WHERE 
            student_id IN (SELECT student_id FROM course_students WHERE course_id = $1)`

	var dbLec []lecture
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, qs, data, &dbLec); err != nil {
		return coursebus.Course{}, fmt.Errorf("db: (fetchStudents) %w", err)
	}

	//======================================================================================================
	// Fetch lectures

	const ql = `
		SELECT 
			title, video_url, public_id, free_preview 
		FROM 
			lectures 
		WHERE 
			id IN (SELECT lecture_id FROM course_curriculum WHERE course_id = $1)`

	var dbStd []student
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, ql, data, &dbStd); err != nil {
		return coursebus.Course{}, fmt.Errorf("db: (fetchStudents) %w", err)
	}

	// Map database records to business logic structures
	course, err := toBusCourse(dbCor)
	if err != nil {
		return coursebus.Course{}, err
	}

	curriculum, err := toBusLectures(dbLec)
	if err != nil {
		return coursebus.Course{}, err
	}

	courseStudents, err := toBusStudents(dbStd)
	if err != nil {
		return coursebus.Course{}, err
	}

	// Attach curriculum and students to the course
	course.Curriculum = curriculum
	course.Students = courseStudents

	return course, nil
}

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
		SELECT completed, completion_date
		FROM CourseProgress
		WHERE user_id = $1 AND course_id = $2`

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
		SELECT lecture_id, viewed, date_viewed
		FROM CourseProgressLectures
		WHERE course_progress_id = (
			SELECT id FROM CourseProgress WHERE user_id = $1 AND course_id = $2
		)`

	var dbLecPro []coursebus.LectureProgress

	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, ql, data, &dbLecPro); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return coursebus.CourseProgress{}, fmt.Errorf("db: %w", coursebus.ErrNotFound)
		}
		return coursebus.CourseProgress{}, fmt.Errorf("db: %w", err)
	}

	//========================================================================================================

	dbCorPro := coursebus.CourseProgress{
		UserID:           userID,
		CourseID:         courseID,
		Completed:        dest.Completed,
		CompletionDate:   dest.CompletionDate,
		LecturesProgress: dbLecPro,
	}

	return dbCorPro, nil
}

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
		UPDATE CourseProgressLectures
			SET viewed = TRUE, date_viewed = $1
		WHERE course_progress_id = (
			SELECT id FROM CourseProgress WHERE user_id = $2 AND course_id = $3
		) AND lecture_id = $4`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

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
		UPDATE CourseProgress
		SET completed = FALSE, completion_date = NULL
		WHERE user_id = $1 AND course_id = $2`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	// Reset lectures progress
	const ql = `
		UPDATE CourseProgressLectures
		SET viewed = FALSE, date_viewed = NULL
		WHERE course_progress_id = (
			SELECT id FROM CourseProgress WHERE user_id = $1 AND course_id = $2
		)`
	if err := sqldb.NamedExecContext(ctx, s.log, s.db, ql, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
