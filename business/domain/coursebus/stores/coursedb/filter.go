package coursedb

import (
	"bytes"
	"strings"

	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
)

func (s *Store) applyFilter(filter coursebus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	// Filter by single ID
	if filter.ID != nil {
		data["course_id"] = *filter.ID
		wc = append(wc, "course_id = :course_id")
	}

	// Filter by multiple categories
	if len(filter.Categories) > 0 {
		data["categories"] = filter.Categories
		wc = append(wc, "category IN (:categories)")
	}

	// Filter by multiple primary languages
	if len(filter.PrimaryLanguages) > 0 {
		data["primary_languages"] = filter.PrimaryLanguages
		wc = append(wc, "primary_language IN (:primary_languages)")
	}

	// Filter by multiple levels
	if len(filter.Levels) > 0 {
		data["levels"] = filter.Levels
		wc = append(wc, "level IN (:levels)")
	}

	// Filter by single price
	if filter.Price != nil {
		data["price"] = *filter.Price
		wc = append(wc, "price = :price")
	}

	// Add WHERE clause if there are conditions
	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}

	// Default ordering and dynamic ordering by fields
	defaultOrder := " ORDER BY price ASC" // Default order is price-lowtohigh
	buf.WriteString(defaultOrder)
}
