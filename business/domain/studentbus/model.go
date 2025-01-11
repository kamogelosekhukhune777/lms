package studentbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/types/name"
)

// StudentCourses represents a student's enrolled courses.
type StudentCourses struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Courses []Course
}

// Course represents an individual course in StudentCourses.
type Course struct {
	CourseID       uuid.UUID
	Title          string
	InstructorID   uuid.UUID
	InstructorName name.Name
	DateOfPurchase time.Time
	CourseImage    string
}
