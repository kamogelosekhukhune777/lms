package coursebus

import (
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/types/name"
)

type Course struct {
	ID              uuid.UUID
	InstructorId    uuid.UUID
	InstructorName  string
	Date            time.Time
	Title           name.Name
	Category        string
	Level           string
	PrimaryLanguage string
	Subtitle        string
	Description     string
	Image           string
	WelcomeMessage  string
	Pricing         float64
	Objectives      string
	IsPublished     bool
}

type NewCourse struct {
	InstructorId    uuid.UUID
	InstructorName  string
	Title           name.Name
	Category        string
	Level           string
	PrimaryLanguage string
	Subtitle        string
	Description     string
	Image           string
	WelcomeMessage  string
	Pricing         float64
}

type UpdateCousre struct {
	Title           *name.Name
	Category        *string
	Level           *string
	PrimaryLanguage *string
	Subtitle        *string
	Description     *string
	Image           *string
	WelcomeMessage  *string
	Pricing         *float64
}
