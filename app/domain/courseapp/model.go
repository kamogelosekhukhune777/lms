package courseapp

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/app/sdk/mid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/types/money"
	"github.com/kamogelosekhukhune777/lms/business/types/name"
)

// Curriculum represents the curriculum of an individual course.
type Curriculum struct {
	Title       string `json:"title"`
	VideoURL    string `json:"videoUrl"`
	FreePreview bool   `json:"freePreview"`
	PublicID    string `json:"public_id"`
}

// Course represents information about an individual course.
type Course struct {
	InstructorID    string       `json:"instructorId"`
	InstructorName  string       `json:"instructorName"`
	Date            string       `json:"date"` //time.Time
	Title           string       `json:"title"`
	Category        string       `json:"category"`
	Level           string       `json:"level"`
	PrimaryLanguage string       `json:"primaryLanguage"`
	Subtitle        string       `json:"subtitle"`
	Description     string       `json:"description"`
	Pricing         float64      `json:"pricing"`
	Objectives      string       `json:"objectives"`
	WelcomeMessage  string       `json:"welcomeMessage"`
	Image           string       `json:"image"`
	Students        []string     `json:"students"`
	Curriculum      []Curriculum `json:"curriculum"`
	IsPublished     bool         `json:"isPublised"`
}

// LectureProgress represents progress on a specific lecture
type LectureProgress struct {
	LectureID  string    `json:"lecture_id"`
	Viewed     bool      `json:"viewed"`
	DateViewed time.Time `json:"date_viewed"`
}

// CourseProgress represents progress on a course for a user
type CourseProgress struct {
	UserID           string            `json:"user_id"`
	CourseID         string            `json:"course_id"`
	Completed        bool              `json:"completed"`
	CompletionDate   time.Time         `json:"completion_date"`
	LecturesProgress []LectureProgress `json:"lectures_progress"`
}

type MarkLectureData struct {
	UserID    string `json:"userId"`
	CourseID  string `json:"courseId"`
	LectureID string `json:"lectureId"`
}
type ResetCourseProgress struct {
	UserID   string `json:"userId"`
	CourseID string `json:"courseId"`
}

func toBusResetCourseProgress(app ResetCourseProgress) (coursebus.ResetCourseProgress, error) {

	userID, err := uuid.Parse(app.UserID)
	if err != nil {
		return coursebus.ResetCourseProgress{}, err
	}

	courseID, err := uuid.Parse(app.CourseID)
	if err != nil {
		return coursebus.ResetCourseProgress{}, err
	}

	bus := coursebus.ResetCourseProgress{
		UserID:   userID,
		CourseID: courseID,
	}

	return bus, nil
}

func toBusMarkLectureData(app MarkLectureData) (coursebus.MarkLectureData, error) {
	userID, err := uuid.Parse(app.UserID)
	if err != nil {
		return coursebus.MarkLectureData{}, err
	}

	courseID, err := uuid.Parse(app.CourseID)
	if err != nil {
		return coursebus.MarkLectureData{}, err
	}

	lectureID, err := uuid.Parse(app.LectureID)
	if err != nil {
		return coursebus.MarkLectureData{}, err
	}

	bus := coursebus.MarkLectureData{
		UserID:    userID,
		CourseID:  courseID,
		LectureID: lectureID,
	}

	return bus, nil
}

func toAppCourse(bus coursebus.Course) Course {
	var curriculum []Curriculum
	var students []string

	if len(bus.Curriculum) == 0 {
		for _, lecture := range bus.Curriculum {
			curriculum = append(curriculum, toAppLecture(lecture))
		}
	}

	if len(bus.Students) == 0 {
		for _, student := range bus.Students {
			studentStr := fmt.Sprintf("ID: %s, CourseID: %s, StudentID: %s, Name: %s, Email: %s, PaidAmount: %s",
				student.ID.String(), student.CourseID.String(), student.StudentID.String(),
				student.StudentName, student.StudentEmail, student.PaidAmount.String())

			students = append(students, studentStr)
		}
	}

	return Course{
		InstructorID:    bus.InstructorID.String(),
		InstructorName:  bus.InstructorName.String(),
		Date:            bus.Date.Format(time.RFC3339),
		Title:           bus.Title.String(),
		Category:        bus.Category,
		Level:           bus.Level,
		PrimaryLanguage: bus.PrimaryLanguage,
		Subtitle:        bus.Subtitle,
		Description:     bus.Description,
		Image:           bus.Image,
		WelcomeMessage:  bus.WelcomeMessage,
		Pricing:         bus.Pricing.Value(),
		Curriculum:      curriculum,
		Students:        students,
		Objectives:      bus.Objectives,
		IsPublished:     bus.IsPublished,
	}
}

func toAppLecture(bus coursebus.Lecture) Curriculum {
	return Curriculum{
		Title:       bus.Title.String(),
		VideoURL:    bus.VideoURL,
		FreePreview: bus.FreePreview,
		PublicID:    bus.PublicID,
	}
}

func toAppCourses(cors []coursebus.Course) []Course {
	app := make([]Course, len(cors))
	for i, cor := range cors {
		app[i] = toAppCourse(cor)
	}

	return app
}

//============================================================================================================

type NewCourse struct {
	InstructorId    string       `json:"instructor_id" validate:"required"`
	InstructorName  string       `json:"instructor_name" validate:"required"`
	Title           string       `json:"title" validate:"required"`
	Pricing         float64      `json:"cost" validate:"required,gte=0"`
	Category        string       `json:"category" validate:"required"`
	Level           string       `json:"level" validate:"required"`
	PrimaryLanguage string       `json:"primary_language" validate:"required"`
	Subtitle        string       `json:"subtitle" validate:"required"`
	Description     string       `json:"description" validate:"required"`
	Image           string       `json:"image" validate:"required"`
	WelcomeMessage  string       `json:"welcome_message" validate:"required"`
	Objectives      string       `json:"objectives" validate:"required"`
	Students        []string     `json:"students"`
	Curriculum      []Curriculum `json:"curriculum" validate:"required"`
	IsPublished     bool         `json:"is_published" validate:"required"`
}

// Validate checks the data in the model is considered clean.
func (app NewCourse) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewCourse(ctx context.Context, app NewCourse) (coursebus.NewCourse, error) {

	instructorID, err := mid.GetInstructorID(ctx)
	if err != nil {
		return coursebus.NewCourse{}, fmt.Errorf("get instructor id: %w", err)
	}

	instructorName, err := name.Parse(app.InstructorName)
	if err != nil {
		return coursebus.NewCourse{}, fmt.Errorf("parse title: %w", err)
	}

	title, err := name.Parse(app.Title)
	if err != nil {
		return coursebus.NewCourse{}, fmt.Errorf("parse title: %w", err)
	}

	price, err := money.Parse(app.Pricing)
	if err != nil {
		return coursebus.NewCourse{}, fmt.Errorf("parse cost: %w", err)
	}

	var curriculum []coursebus.Lecture
	curriculum = make([]coursebus.Lecture, 0, len(app.Curriculum))

	for _, lecture := range app.Curriculum {
		if lec, err := toBusCurriculum(lecture); err == nil {
			curriculum = append(curriculum, lec)
		} else {
			return coursebus.NewCourse{}, fmt.Errorf("invalid lecture: %w", err)
		}
	}

	bus := coursebus.NewCourse{
		InstructorID:    instructorID,
		InstructorName:  instructorName,
		Title:           title,
		Pricing:         price,
		Category:        app.Category,
		Level:           app.Level,
		PrimaryLanguage: app.PrimaryLanguage,
		Subtitle:        app.Subtitle,
		Description:     app.Description,
		Image:           app.Image,
		WelcomeMessage:  app.WelcomeMessage,
		Curriculum:      curriculum,
		Students:        app.Students,
		Objectives:      app.Objectives,
	}

	return bus, nil
}

func toBusCurriculum(app Curriculum) (coursebus.Lecture, error) {
	title, err := name.Parse(app.Title)
	if err != nil {
		return coursebus.Lecture{}, fmt.Errorf("parse title: %w", err)
	}

	return coursebus.Lecture{
		Title:       title,
		VideoURL:    app.VideoURL,
		FreePreview: app.FreePreview,
		PublicID:    app.PublicID,
	}, nil
}

//============================================================================================================

type UpdateCourse struct {
	Title           *string  `json:"title" validate:"required"`
	Pricing         *float64 `json:"cost" validate:"required,gte=0"`
	Category        *string  `json:"category" validate:"required"`
	Level           *string  `json:"level" validate:"required"`
	PrimaryLanguage *string  `json:"primary_language" validate:"required"`
	Subtitle        *string  `json:"subtitle" validate:"required"`
	Description     *string  `json:"description" validate:"required"`
	Image           *string  `json:"image" validate:"required"`
	WelcomeMessage  *string  `json:"welcome_message" validate:"required"`
	Objectives      *string  `json:"objectives" validate:"required"`
	IsPublished     *bool    `json:"is_published" validate:"required"`
}

func toBusUpdateCourse(app UpdateCourse) (coursebus.UpdateCourse, error) {

	var title *name.Name
	if app.Title != nil {
		nm, err := name.Parse(*app.Title)
		if err != nil {
			return coursebus.UpdateCourse{}, fmt.Errorf("parse: %w", err)
		}
		title = &nm
	}

	var pricing *money.Money
	if app.Pricing != nil {
		price, err := money.Parse(*app.Pricing)
		if err != nil {
			return coursebus.UpdateCourse{}, fmt.Errorf("parse: %w", err)
		}
		pricing = &price
	}

	bus := coursebus.UpdateCourse{
		Title:           title,
		Pricing:         pricing,
		Category:        app.Category,
		Level:           app.Level,
		PrimaryLanguage: app.PrimaryLanguage,
		Subtitle:        app.Subtitle,
		Description:     app.Description,
		Image:           app.Image,
		WelcomeMessage:  app.WelcomeMessage,
	}

	return bus, nil
}
