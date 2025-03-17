package coursedb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/types/money"
)

type course struct {
	ID              uuid.UUID `db:"course_id"`
	InstructorID    uuid.UUID `db:"instructor_id"`
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
	CreatedAt       time.Time `db:"created_at"`
}

func toDBCourse(bus coursebus.Course) course {
	return course{
		ID:              bus.ID,
		InstructorID:    bus.InstructorID,
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
		CreatedAt:       bus.CreatedAt.UTC(),
	}
}

func toBusCourse(db course) (coursebus.Course, error) {

	price, err := money.Parse(db.Pricing)
	if err != nil {
		return coursebus.Course{}, fmt.Errorf("parse cost: %w", err)
	}

	bus := coursebus.Course{
		ID:              db.ID,
		InstructorID:    db.ID,
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
		CreatedAt:       db.CreatedAt.In(time.Local),
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

//=============================================================================================================================

type lecture struct {
	ID          uuid.UUID
	CourseID    uuid.UUID
	Title       string
	VideoURL    string
	PublicID    string
	FreePreview bool
}

func toDBLecture(bus coursebus.Lecture) lecture {
	return lecture{
		ID:          bus.ID,
		CourseID:    bus.CourseID,
		Title:       bus.Title,
		VideoURL:    bus.VideoURL,
		PublicID:    bus.PublicID,
		FreePreview: bus.FreePreview,
	}
}

func toBusLecture(db lecture) (coursebus.Lecture, error) {

	bus := coursebus.Lecture{
		ID:          db.ID,
		CourseID:    db.CourseID,
		Title:       db.Title,
		VideoURL:    db.VideoURL,
		PublicID:    db.PublicID,
		FreePreview: db.FreePreview,
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

type student struct {
	ID         uuid.UUID
	StudentID  uuid.UUID
	CourseID   uuid.UUID
	PaidAmount float64
	EnrolledAt time.Time
}

func toDBStudent(bus coursebus.Student) student {

	return student{
		ID:         bus.ID,
		StudentID:  bus.StudentID,
		CourseID:   bus.CourseID,
		PaidAmount: bus.PaidAmount.Value(),
		EnrolledAt: bus.EnrolledAt,
	}
}

func toBusStudent(db student) (coursebus.Student, error) {
	paid, err := money.Parse(db.PaidAmount)
	if err != nil {
		return coursebus.Student{}, fmt.Errorf("parse cost: %w", err)
	}

	bus := coursebus.Student{
		ID:         db.ID,
		StudentID:  db.StudentID,
		CourseID:   db.CourseID,
		PaidAmount: paid,
		EnrolledAt: db.EnrolledAt.In(time.Local),
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
