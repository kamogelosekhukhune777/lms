package coursedb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/types/name"
)

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

func toDBCourse(bus coursebus.Course) course {
	db := course{
		ID:              bus.ID,
		InstructorId:    bus.InstructorId,
		InstructorName:  bus.InstructorName,
		Date:            bus.Date.UTC(),
		Title:           bus.Title.String(),
		Category:        bus.Category,
		Level:           bus.Level,
		PrimaryLanguage: bus.PrimaryLanguage,
		Subtitle:        bus.Subtitle,
		Description:     bus.Description,
		Image:           bus.Image,
		WelcomeMessage:  bus.WelcomeMessage,
		Pricing:         bus.Pricing,
		Objectives:      bus.Objectives,
		IsPublished:     bus.IsPublished,
	}

	return db
}

func toBusCourse(db course) (coursebus.Course, error) {

	title, err := name.Parse(db.Title)
	if err != nil {
		return coursebus.Course{}, fmt.Errorf("parse title: %w", err)
	}

	bus := coursebus.Course{
		ID:              db.ID,
		InstructorId:    db.InstructorId,
		InstructorName:  db.InstructorName,
		Date:            db.Date,
		Title:           title,
		Category:        db.Category,
		Level:           db.Level,
		PrimaryLanguage: db.PrimaryLanguage,
		Subtitle:        db.Subtitle,
		Description:     db.Description,
		Image:           db.Image,
		WelcomeMessage:  db.WelcomeMessage,
		Pricing:         db.Pricing,
		Objectives:      db.Objectives,
		IsPublished:     db.IsPublished,
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
