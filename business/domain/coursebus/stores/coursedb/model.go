package coursedb

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/types/money"
	"github.com/kamogelosekhukhune777/lms/business/types/name"
)

// LectureProgress represents progress for a specific lecture.
type lectureProgress struct {
	LectureID  uuid.UUID `db:"lecture_id"`
	Viewed     bool      `db:"viewed"`
	DateViewed time.Time `db:"date_viewed"`
}

// CourseProgress represents a student's progress for a course.
type courseProgress struct {
	ID               uuid.UUID         `db:"id"`
	UserID           uuid.UUID         `db:"user_id"`
	CourseID         uuid.UUID         `db:"course_id"`
	Completed        bool              `db:"completed"`
	CompletionDate   time.Time         `db:"completion_date"`
	LecturesProgress []lectureProgress `db:"lectures_progress"` // Handle serialization for storing JSON-like data in SQL if needed
}

// Lecture represents an individual lecture in a course.
type lecture struct {
	ID          uuid.UUID `db:"lecture_id"`
	Title       string    `db:"title"`
	VideoURL    string    `db:"video_url"`
	PublicID    string    `db:"public_id"`
	FreePreview bool      `db:"free_preview"`
}

// CourseSchema represents the schema for a course.
type courseSchema struct {
	ID              uuid.UUID `db:"id"`
	InstructorID    uuid.UUID `db:"instructor_id"`
	InstructorName  string    `db:"instructor_name"`
	Date            time.Time `db:"date"`
	Title           string    `db:"title"`
	Category        string    `db:"category"`
	Level           string    `db:"level"`
	PrimaryLanguage string    `db:"primary_language"`
	Subtitle        string    `db:"subtitle"`
	Description     string    `db:"description"`
	Image           string    `db:"image"`
	WelcomeMessage  string    `db:"welcome_message"`
	Pricing         float64   `db:"pricing"`
	Objectives      string    `db:"objectives"`
	Students        []student `db:"students"`
	Curriculum      []lecture `db:"curriculum"`
	IsPublished     bool      `db:"is_published"`
}

// Student represents a student enrolled in a course.
type student struct {
	ID           uuid.UUID `db:"course_student_id"`
	CourseID     uuid.UUID `db:"course_id"`
	StudentID    uuid.UUID `db:"student_id"`
	StudentName  string    `db:"student_name"`
	StudentEmail string    `db:"student_email"`
	PaidAmount   float64   `db:"paid_amount"`
}

//======================================================================================================================
//toDB(database)

func toDBCourseSchema(bus coursebus.CourseSchema) courseSchema {
	return courseSchema{
		ID:              bus.ID,
		InstructorID:    bus.InstructorID,
		InstructorName:  bus.InstructorName.String(),
		Date:            bus.Date,
		Title:           bus.Title,
		Category:        bus.Category,
		Level:           bus.Level,
		PrimaryLanguage: bus.PrimaryLanguage,
		Subtitle:        bus.Subtitle,
		Description:     bus.Description,
		Image:           bus.Image,
		WelcomeMessage:  bus.WelcomeMessage,
		Pricing:         bus.Pricing.Value(),
		Objectives:      bus.Objectives,
		IsPublished:     bus.IsPublished,
	}
}

func toDBStudent(db coursebus.Student) student {
	return student{
		StudentID:    db.StudentID,
		StudentName:  db.StudentName.String(),
		StudentEmail: db.StudentEmail.String(),
		PaidAmount:   db.PaidAmount.Value(),
	}

}

func toDBLecture(db coursebus.Lecture) lecture {

	return lecture{
		Title:       db.Title,
		VideoURL:    db.VideoURL,
		PublicID:    db.PublicID,
		FreePreview: db.FreePreview,
	}
}

// ======================================================================================================================
// toBusiness from db

func toBusCourseSchema(db courseSchema) (coursebus.CourseSchema, error) {
	name, err := name.Parse(db.InstructorName)
	if err != nil {
		return coursebus.CourseSchema{}, fmt.Errorf("parse: %w", err)
	}

	price, err := money.Parse(db.Pricing)
	if err != nil {
		return coursebus.CourseSchema{}, fmt.Errorf("parse: %w", err)
	}

	//students and lectures

	bus := coursebus.CourseSchema{
		ID:              db.ID,
		InstructorID:    db.ID,
		InstructorName:  name,
		Date:            db.Date,
		Title:           db.Title,
		Category:        db.Category,
		Level:           db.Level,
		PrimaryLanguage: db.PrimaryLanguage,
		Subtitle:        db.Subtitle,
		Description:     db.Description,
		Image:           db.Image,
		WelcomeMessage:  db.WelcomeMessage,
		Pricing:         price,
		Objectives:      db.Objectives,
		IsPublished:     db.IsPublished,
	}

	return bus, nil
}

func toBusLecture(db lecture) (coursebus.Lecture, error) {
	bus := coursebus.Lecture{
		Title:       db.Title,
		VideoURL:    db.VideoURL,
		PublicID:    db.PublicID,
		FreePreview: db.FreePreview,
	}

	return bus, nil
}

func toBusStudent(db student) (coursebus.Student, error) {
	name, err := name.Parse(db.StudentName)
	if err != nil {
		return coursebus.Student{}, nil
	}

	email := mail.Address{
		Address: db.StudentEmail,
	}

	money, err := money.Parse(db.PaidAmount)
	if err != nil {
		return coursebus.Student{}, nil
	}

	bus := coursebus.Student{
		StudentID:    db.StudentID,
		StudentName:  name,
		StudentEmail: email,
		PaidAmount:   money,
	}

	return bus, nil
}

//======================================================================================================================
//toBusiness(slices)

func toBusCoursesSchema(dbs []courseSchema) ([]coursebus.CourseSchema, error) {
	bus := make([]coursebus.CourseSchema, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusCourseSchema(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toBusLectures(dbs []lecture) ([]coursebus.Lecture, error) {
	bus := make([]coursebus.Lecture, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusLecture(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toBusStudents(dbs []student) ([]coursebus.Student, error) {
	bus := make([]coursebus.Student, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusStudent(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

// ()

func toBusLectureProgress(dbs []lectureProgress) []coursebus.LectureProgress {
	bus := make([]coursebus.LectureProgress, len(dbs))

	for _, db := range dbs {
		bus = append(bus, coursebus.LectureProgress{
			LectureID:  db.LectureID,
			Viewed:     db.Viewed,
			DateViewed: db.DateViewed,
		})
	}

	return bus
}
