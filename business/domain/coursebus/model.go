package coursebus

import (
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/types/money"
)

// Course represents an individual course.
type Course struct {
	ID              uuid.UUID
	InstructorID    uuid.UUID
	Title           string
	Category        string
	Level           string
	PrimaryLanguage string
	Subtitle        string
	Description     string
	Image           string
	WelcomeMessage  string
	Pricing         money.Money
	Objectives      string
	IsPublished     bool
	CreatedAt       time.Time
}

// NewCourse is what we require from clients when adding a Course.
type NewCourse struct {
	InstructorID    uuid.UUID
	Title           string
	Category        string
	Level           string
	PrimaryLanguage string
	Subtitle        string
	Description     string
	Image           string
	WelcomeMessage  string
	Pricing         money.Money
	Objectives      string
}

type UpdateCourse struct {
	Title           *string
	Category        *string
	Level           *string
	PrimaryLanguage *string
	Subtitle        *string
	Description     *string
	Image           *string
	WelcomeMessage  *string
	Pricing         *money.Money
	Objectives      *string
}

//=================================================================================
