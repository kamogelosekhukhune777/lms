package orderbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/types/money"
)

type Order struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	UserName       string
	UserEmail      string
	OrderStatus    string
	PaymentMethod  string
	PaymentStatus  string
	OrderDate      time.Time
	PaymentID      string
	PayerID        string
	InstructorID   uuid.UUID
	InstructorName string
	CourseImage    string
	CourseTitle    string
	CourseID       uuid.UUID
	CoursePricing  money.Money
}

type NewOrder struct {
	UserID         uuid.UUID
	UserName       string
	UserEmail      string
	OrderStatus    string
	PaymentMethod  string
	PaymentStatus  string
	OrderDate      time.Time
	PaymentID      string
	PayerID        string
	InstructorID   uuid.UUID
	InstructorName string
	CourseImage    string
	CourseTitle    string
	CourseID       uuid.UUID
	CoursePricing  money.Money
}
