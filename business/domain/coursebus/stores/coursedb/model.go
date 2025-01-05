package coursedb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/types/money"
	"github.com/kamogelosekhukhune777/lms/business/types/name"
)

type lecture struct {
	ID          uuid.UUID `json:"lecture_id"`
	Title       string    `json:"title"`
	VideoURL    string    `json:"video_url"`
	PublicID    string    `json:"public_id"`
	FreePreview bool      `json:"free_preview"`
}

// Student represents the schema for a student in a course.
type student struct {
	ID           uuid.UUID `db:"course_student_id"`
	CourseID     uuid.UUID `db:"course_id"`
	StudentID    uuid.UUID `db:"student_id"`
	StudentName  string    `db:"student_name"`
	StudentEmail string    `db:"student_email"`
	PaidAmount   float64   `db:"paid_amount"`
}

type course struct {
	ID              uuid.UUID `db:"course_id"`
	InstructorId    uuid.UUID `db:"instructor_id"`
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
	IsPublished     bool      `db:"is_published"`
}

type LectureProgress struct {
	LectureID  string    `db:"lecture_id"`
	Viewed     bool      `db:"viewed"`
	DateViewed time.Time `db:"date_viewed"`
}

type CourseProgress struct {
	UserID           string            `db:"user_id"`
	CourseID         string            `db:"course_id"`
	Completed        bool              `db:"completed"`
	CompletionDate   time.Time         `db:"completion_date"`
	LecturesProgress []LectureProgress `db:"lectures_progress"`
}

type curriculum struct {
	CourseID  uuid.UUID `db:"course_id"`
	LectureID uuid.UUID `db:"lecture_id"`
}

func toDBCurriculum(lectureID, courseID uuid.UUID) curriculum {
	db := curriculum{
		CourseID:  courseID,
		LectureID: lectureID,
	}

	return db
}

func toDBCourse(bus coursebus.Course) course {
	db := course{
		ID:              bus.ID,
		InstructorId:    bus.InstructorID,
		InstructorName:  bus.InstructorName.String(),
		Date:            bus.Date.UTC(),
		Title:           bus.Title.String(),
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

	return db
}

func toDBLecture(bus coursebus.Lecture) lecture {
	db := lecture{
		ID:          bus.ID,
		Title:       bus.Title.String(),
		VideoURL:    bus.VideoURL,
		FreePreview: bus.FreePreview,
		PublicID:    bus.PublicID,
	}

	return db
}

func toDBStudent(bus coursebus.Student) student {
	db := student{
		ID:           bus.ID,
		CourseID:     bus.CourseID,
		StudentID:    bus.StudentID,
		StudentEmail: bus.StudentEmail,
		PaidAmount:   bus.PaidAmount.Value(),
	}

	return db
}

func toBusCourse(db course) (coursebus.Course, error) {

	title, err := name.Parse(db.Title)
	if err != nil {
		return coursebus.Course{}, fmt.Errorf("parse title: %w", err)
	}

	instructorName, err := name.Parse(db.InstructorName)
	if err != nil {
		return coursebus.Course{}, fmt.Errorf("parse title: %w", err)
	}

	price, err := money.Parse(db.Pricing)
	if err != nil {
		return coursebus.Course{}, fmt.Errorf("parse pricing: %w", err)
	}

	bus := coursebus.Course{
		ID:              db.ID,
		InstructorID:    db.InstructorId,
		InstructorName:  instructorName,
		Date:            db.Date,
		Title:           title,
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
	title, err := name.Parse(db.Title)
	if err != nil {
		return coursebus.Lecture{}, fmt.Errorf("parse title: %w", err)
	}

	bus := coursebus.Lecture{
		ID:          db.ID,
		Title:       title,
		VideoURL:    db.VideoURL,
		FreePreview: db.FreePreview,
		PublicID:    db.PublicID,
	}

	return bus, nil
}

func toBusStudent(db student) (coursebus.Student, error) {
	studentName, err := name.Parse(db.StudentName)
	if err != nil {
		return coursebus.Student{}, fmt.Errorf("parse title: %w", err)
	}

	paidAmount, err := money.Parse(db.PaidAmount)
	if err != nil {
		return coursebus.Student{}, fmt.Errorf("parse pricing: %w", err)
	}

	bus := coursebus.Student{
		ID:           db.ID,
		CourseID:     db.CourseID,
		StudentID:    db.StudentID,
		StudentName:  studentName,
		StudentEmail: db.StudentEmail,
		PaidAmount:   paidAmount,
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

func toBusLectures(dbls []lecture) ([]coursebus.Lecture, error) {
	bus := make([]coursebus.Lecture, len(dbls))

	for i, dbl := range dbls {
		var err error
		bus[i], err = toBusLecture(dbl)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func toBusCourses(dbs []course) ([]coursebus.Course, error) {
	bus := make([]coursebus.Course, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusCourse(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

//===============================================================================
//===============================================================================
