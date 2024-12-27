// Package mid provides app level middleware support.
package mid

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/sqldb"
)

type ctxKey int

const (
	claimKey ctxKey = iota + 1
	InstructorIDKey
	courseKey
	trKey
)

func setInstructorID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, InstructorIDKey, userID)
}

// GetInstructorID returns the user id from the context.
func GetInstructorID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(InstructorIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("instructor id not found in context")
	}

	return v, nil
}

func setCourse(ctx context.Context, cor coursebus.Course) context.Context {
	return context.WithValue(ctx, courseKey, cor)
}

// GetCourse returns the user from the context.
func GetCourse(ctx context.Context) (coursebus.Course, error) {
	v, ok := ctx.Value(courseKey).(coursebus.Course)
	if !ok {
		return coursebus.Course{}, errors.New("cousre not found in context")
	}

	return v, nil
}

func setTran(ctx context.Context, tx sqldb.CommitRollbacker) context.Context {
	return context.WithValue(ctx, trKey, tx)
}

// GetTran retrieves the value that can manage a transaction.
func GetTran(ctx context.Context) (sqldb.CommitRollbacker, error) {
	v, ok := ctx.Value(trKey).(sqldb.CommitRollbacker)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}

	return v, nil
}
