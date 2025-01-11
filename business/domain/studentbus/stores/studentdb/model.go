package studentdb

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/domain/studentbus"
	"github.com/kamogelosekhukhune777/lms/business/types/money"
	"github.com/kamogelosekhukhune777/lms/business/types/name"
)

// StudentCourses represents a student's enrolled courses.
type studentCourses struct {
	ID      uuid.UUID `db:"id"`
	UserID  uuid.UUID `db:"user_id"`
	Courses []course  `db:"courses"` // Handle serialization for storing JSON-like data in SQL if needed
}

// Course represents an individual course in StudentCourses.
type course struct {
	CourseID       uuid.UUID `db:"course_id"`
	Title          string    `db:"title"`
	InstructorID   uuid.UUID `db:"instructor_id"`
	InstructorName string    `db:"instructor_name"`
	DateOfPurchase time.Time `db:"date_of_purchase"`
	CourseImage    string    `db:"course_image"`
}

// LectureProgress represents progress for a specific lecture.
type lectureProgress struct {
	LectureID  uuid.UUID `db:"lecture_id"`
	Viewed     bool      `db:"viewed"`
	DateViewed time.Time `db:"date_viewed"`
}

// CourseProgress represents a student's progress for a course.
type courseProgress struct {
	ID               uuid.UUID         `db:"id"`
	UserID           uuid.UUID         `db:"user_id"`
	CourseID         uuid.UUID         `db:"course_id"`
	Completed        bool              `db:"completed"`
	CompletionDate   time.Time         `db:"completion_date"`
	LecturesProgress []lectureProgress `db:"lectures_progress"` // Handle serialization for storing JSON-like data in SQL if needed
}

// Lecture represents an individual lecture in a course.
type lecture struct {
	ID          uuid.UUID `db:"lecture_id"`
	Title       string    `db:"title"`
	VideoURL    string    `db:"video_url"`
	PublicID    string    `db:"public_id"`
	FreePreview bool      `db:"free_preview"`
}

// Student represents a student enrolled in a course.
type student struct {
	ID           uuid.UUID `db:"course_student_id"`
	CourseID     uuid.UUID `db:"course_id"`
	StudentID    uuid.UUID `db:"student_id"`
	StudentName  string    `db:"student_name"`
	StudentEmail string    `db:"student_email"`
	PaidAmount   float64   `db:"paid_amount"`
}

//======================================================================================================================
//toDB(database)

func toDBStudent(db coursebus.Student) student {
	return student{
		StudentID:    db.StudentID,
		StudentName:  db.StudentName.String(),
		StudentEmail: db.StudentEmail.String(),
		PaidAmount:   db.PaidAmount.Value(),
	}

}

func toDBLecture(db coursebus.Lecture) lecture {

	return lecture{
		Title:       db.Title,
		VideoURL:    db.VideoURL,
		PublicID:    db.PublicID,
		FreePreview: db.FreePreview,
	}
}

// ======================================================================================================================
// toBusiness from db

func toBusLecture(db lecture) (coursebus.Lecture, error) {
	bus := coursebus.Lecture{
		Title:       db.Title,
		VideoURL:    db.VideoURL,
		PublicID:    db.PublicID,
		FreePreview: db.FreePreview,
	}

	return bus, nil
}

func toBusStudent(db student) (coursebus.Student, error) {
	name, err := name.Parse(db.StudentName)
	if err != nil {
		return coursebus.Student{}, nil
	}

	email := mail.Address{
		Address: db.StudentEmail,
	}

	money, err := money.Parse(db.PaidAmount)
	if err != nil {
		return coursebus.Student{}, nil
	}

	bus := coursebus.Student{
		StudentID:    db.StudentID,
		StudentName:  name,
		StudentEmail: email,
		PaidAmount:   money,
	}

	return bus, nil
}

//======================================================================================================================
//toBusiness(slices)

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

func toBusLectureProgress(dbs []lectureProgress) []coursebus.LectureProgress {
	bus := make([]coursebus.LectureProgress, len(dbs))

	for _, db := range dbs {
		bus = append(bus, coursebus.LectureProgress{
			LectureID:  db.LectureID,
			Viewed:     db.Viewed,
			DateViewed: db.DateViewed,
		})
	}

	return bus
}

func toBusStudentCourses(dbs []course) ([]studentbus.Course, error) {
	bus := make([]studentbus.Course, len(dbs))

	for _, db := range dbs {
		name, err := name.Parse(db.InstructorName)
		if err != nil {
			return []studentbus.Course{}, nil
		}

		bus = append(bus, studentbus.Course{
			CourseID:       db.CourseID,
			Title:          db.Title,
			InstructorID:   db.InstructorID,
			InstructorName: name,
			DateOfPurchase: db.DateOfPurchase,
			CourseImage:    db.CourseImage,
		})
	}

	return bus, nil
}
