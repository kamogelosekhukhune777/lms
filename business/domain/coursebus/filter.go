package coursebus

import (
	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/business/types/money"
)

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type QueryFilter struct {
	ID               *uuid.UUID
	Categories       []string // Allow multiple categories
	PrimaryLanguages []string // Allow multiple primary languages
	Levels           []string // Allow multiple levels
	Price            *money.Money
}
