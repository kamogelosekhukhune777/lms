package coursebus

import (
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/types/money"
	"github.com/kamogelosekhukhune777/lms/business/types/name"
)

type Course struct {
	ID              uuid.UUID
	InstructorID    uuid.UUID
	InstructorName  name.Name
	Date            time.Time
	Title           name.Name
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

type Lecture struct {
	ID          uuid.UUID
	Title       name.Name
	VideoURL    string
	FreePreview bool
	PublicID    string
}

type Student struct {
	ID           uuid.UUID
	CourseID     uuid.UUID
	StudentID    uuid.UUID
	StudentName  name.Name
	StudentEmail string
	PaidAmount   money.Money
}

//=================================================================================

type NewCourse struct {
	InstructorID    uuid.UUID
	InstructorName  name.Name
	Title           name.Name
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
	StudentEmail string
	PaidAmount   money.Money
}

//=================================================================================

type UpdateCourse struct {
	Title           *name.Name
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
	Title       *name.Name
	VideoURL    *string
	FreePreview *bool
	PublicID    *string
}

type UpdateStudent struct {
	CourseID     *uuid.UUID
	StudentID    *uuid.UUID
	StudentName  *name.Name
	StudentEmail *string
	PaidAmount   *money.Money
}
