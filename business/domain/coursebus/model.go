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
	Curriculum      []Lecture
	Student         []Student
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
	Curriculum      []Lecture
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
	Curriculum      []Lecture
	Pricing         *money.Money
	Objectives      *string
}

//========================================================================================================================
//========================================================================================================================

type Lecture struct {
	ID          uuid.UUID
	CourseID    uuid.UUID
	Title       string
	VideoURL    string
	PublicID    string
	FreePreview bool
}

// Student(course Students)/(Enrollments)
type Student struct {
	ID         uuid.UUID
	StudentID  uuid.UUID
	CourseID   uuid.UUID
	PaidAmount money.Money
	EnrolledAt time.Time
}

//========================================================================================================================
//========================================================================================================================

type CourseProgress struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	CourseID       uuid.UUID
	Completed      bool
	CompletionDate time.Time
}

type LectureProgress struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	LectureID  uuid.UUID
	Viewed     bool
	DateViewed time.Time
}
