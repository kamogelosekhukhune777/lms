package courseapp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/app/sdk/mid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/types/money"
)

// Course represents information about an individual course.
type Course struct {
	ID              string    `json:"course_id"`
	InstructorID    string    `json:"instructor_id"`
	Title           string    `json:"title"`
	Category        string    `json:"category"`
	Level           string    `json:"level"`
	PrimaryLanguage string    `json:"primary_language"`
	Subtitle        string    `json:"subtitle"`
	Description     string    `json:"description"`
	Image           string    `json:"image"`
	WelcomeMessage  string    `json:"welcome_message"`
	Pricing         float64   `json:"pricing"`
	Objectives      string    `json:"objectives"`
	Curriculum      []Lecture `json:"curriculum"`
	IsPublished     bool      `json:"is_published"`
	CreatedAt       time.Time `json:"created_at"`
}

// Encode implements the encoder interface.
func (app Course) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppCourse(cor coursebus.Course) Course {
	return Course{
		ID:              cor.ID.String(),
		InstructorID:    cor.InstructorID.String(),
		Title:           cor.Title,
		Category:        cor.Category,
		Level:           cor.Level,
		PrimaryLanguage: cor.PrimaryLanguage,
		Subtitle:        cor.Subtitle,
		Description:     cor.Description,
		Image:           cor.Image,
		WelcomeMessage:  cor.WelcomeMessage,
		Pricing:         cor.Pricing.Value(),
		Objectives:      cor.Objectives,
		IsPublished:     cor.IsPublished,
		CreatedAt:       cor.CreatedAt.In(time.Local),
	}
}

func toAppCourses(cors []coursebus.Course) Courses {
	app := make([]Course, len(cors))
	for i, cor := range cors {
		app[i] = toAppCourse(cor)
	}

	return app
}

type Courses []Course

func (app Courses) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

type BoolResult bool

func (app BoolResult) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// =============================================================================

// NewCourse defines the data needed to add a new course.
type NewCourse struct {
	Title           string  `json:"title" validate:"required"`
	Category        string  `json:"category" validate:"required"`
	Level           string  `json:"level" validate:"required"`
	PrimaryLanguage string  `json:"primary_language" validate:"required"`
	Subtitle        string  `json:"subtitle" validate:"required"`
	Description     string  `json:"description" validate:"required"`
	Image           string  `json:"image" validate:"required"`
	WelcomeMessage  string  `json:"welcome_message" validate:"required"`
	Pricing         float64 `json:"pricing" validate:"required,gte=0"`
	Objectives      string  `json:"objectives" validate:"required"`
	IsPublished     bool    `json:"is_published" validate:"required"`
}

// Decode implements the decoder interface.
func (app *NewCourse) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewCourse) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewCourse(ctx context.Context, app NewCourse) (coursebus.NewCourse, error) {
	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return coursebus.NewCourse{}, fmt.Errorf("getuserid: %w", err)
	}

	price, err := money.Parse(app.Pricing)
	if err != nil {
		return coursebus.NewCourse{}, fmt.Errorf("parse cost: %w", err)
	}

	bus := coursebus.NewCourse{
		InstructorID:    userID,
		Title:           app.Title,
		Category:        app.Category,
		Level:           app.Level,
		PrimaryLanguage: app.PrimaryLanguage,
		Subtitle:        app.Subtitle,
		Description:     app.Description,
		Image:           app.Image,
		WelcomeMessage:  app.WelcomeMessage,
		Pricing:         price,
		Objectives:      app.Objectives,
	}

	return bus, nil
}

// =============================================================================

// UpdateProduct defines the data needed to update a product.
type UpdateCourse struct {
	Title           *string  `json:"title"`
	Category        *string  `json:"category"`
	Level           *string  `json:"level"`
	PrimaryLanguage *string  `josn:"primary_language"`
	Subtitle        *string  `json:"subtitle"`
	Description     *string  `josn:"description"`
	Image           *string  `json:"image"`
	WelcomeMessage  *string  `json:"welcome_message"`
	Pricing         *float64 `json:"pricing" validate:"omitempty,gte=0"`
	Objectives      *string  `json:"objectives"`
}

// Decode implements the decoder interface.
func (app *UpdateCourse) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app UpdateCourse) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusUpdateCourse(app UpdateCourse) (coursebus.UpdateCourse, error) {

	var price *money.Money
	if app.Pricing != nil {
		prc, err := money.Parse(*app.Pricing)
		if err != nil {
			return coursebus.UpdateCourse{}, fmt.Errorf("parse: %w", err)
		}
		price = &prc
	}

	bus := coursebus.UpdateCourse{
		Title:           app.Title,
		Category:        app.Category,
		Level:           app.Level,
		PrimaryLanguage: app.PrimaryLanguage,
		Subtitle:        app.Subtitle,
		Description:     app.Description,
		Image:           app.Image,
		WelcomeMessage:  app.WelcomeMessage,
		Pricing:         price,
		Objectives:      app.Objectives,
	}

	return bus, nil
}

//=============================================================

type Lecture struct {
	ID          string `json:"lecture_id"`
	CourseID    string `json:"course_id"`
	Title       string `json:"title"`
	VideoURL    string `json:"video_url"`
	PublicID    string `json:"public_id"`
	FreePreview bool   `json:"free_preview"`
}

// Encode implements the encoder interface.
func (app Lecture) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

//=====================================================================

type NewLecture struct{}

// Decode implements the decoder interface.
func (app *NewLecture) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewLecture) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

//========================================================================

// Student(course Students)/(Enrollments)
type Student struct {
	ID         string    `json:"id"`
	StudentID  string    `json:"student_id"`
	CourseID   string    `json:"course_id"`
	PaidAmount float64   `json:"paid_amount"`
	EnrolledAt time.Time `json:"enrolled_at"`
}

// Encode implements the encoder interface.
func (app Student) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

//===========================================================================

type NewStudent struct{}

// Decode implements the decoder interface.
func (app *NewStudent) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewStudent) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

//==============================================================================

type CourseProgess struct {
	ID             string
	UserID         string
	CourseID       string
	Completed      bool
	CompletionDate time.Time
}

// Encode implements the encoder interface.
func (app CourseProgess) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppCourseProgress(corp coursebus.CourseProgress) CourseProgess {
	return CourseProgess{
		ID:             corp.ID.String(),
		UserID:         corp.UserID.String(),
		CourseID:       corp.CourseID.String(),
		Completed:      corp.Completed,
		CompletionDate: corp.CompletionDate,
	}
}
