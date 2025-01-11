package coursebus

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/types/money"
	"github.com/kamogelosekhukhune777/lms/business/types/name"
)

// LectureProgress represents progress for a specific lecture.
type LectureProgress struct {
	LectureID  uuid.UUID
	Viewed     bool
	DateViewed time.Time
}

// CourseProgress represents a student's progress for a course.
type CourseProgress struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	CourseID         uuid.UUID
	Completed        bool
	CompletionDate   time.Time
	LecturesProgress []LectureProgress
}

// Lecture represents an individual lecture in a course.
type Lecture struct {
	Title       string
	VideoURL    string
	PublicID    string
	FreePreview bool
}

// CourseSchema represents the schema for a course.
type CourseSchema struct {
	ID              uuid.UUID
	InstructorID    uuid.UUID
	InstructorName  name.Name
	Date            time.Time
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
	Students        []Student
	Curriculum      []Lecture
	IsPublished     bool
}

// Student represents a student enrolled in a course.
type Student struct {
	StudentID    uuid.UUID
	StudentName  name.Name
	StudentEmail mail.Address
	PaidAmount   money.Money
}

//===

type MarkLectureData struct {
	UserID    uuid.UUID
	CourseID  uuid.UUID
	LectureID uuid.UUID
}
type ResetCourseProgresParams struct {
	UserID   uuid.UUID
	CourseID uuid.UUID
}

//=================================================================================

type NewCourseSchema struct {
	InstructorID    uuid.UUID
	InstructorName  name.Name
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
	Students        []string
	Curriculum      []Lecture
	IsPublished     bool
}

type NewLecture struct {
	Title       name.Name
	VideoURL    string
	FreePreview bool
	PublicID    string
}

type NewStudent struct {
	CourseID     uuid.UUID
	StudentID    uuid.UUID
	StudentName  name.Name
	StudentEmail mail.Address
	PaidAmount   money.Money
}

//=================================================================================

type UpdateCourseSchema struct {
	Title           *string
	Category        *string
	Level           *string
	PrimaryLanguage *string
	Subtitle        *string
	Description     *string
	Image           *string
	WelcomeMessage  *string
	Pricing         *money.Money
	Students        []Student
	Curriculum      []Lecture
}

type UpdateLecture struct {
	Title       *string
	VideoURL    *string
	FreePreview *bool
	PublicID    *string
}

type UpdateStudent struct {
	CourseID     *uuid.UUID
	StudentID    *uuid.UUID
	StudentName  *name.Name
	StudentEmail *mail.Address
	PaidAmount   *money.Money
}
