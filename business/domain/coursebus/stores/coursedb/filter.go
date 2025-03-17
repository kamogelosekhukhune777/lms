package coursedb

import (
	"bytes"
	"strings"

	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
)

func (s *Store) applyFilter(filter coursebus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.Category != nil {
		data["category"] = *filter.Category
		wc = append(wc, "category = :category")
	}

	if filter.Level != nil {
		data["level"] = *filter.Level
		wc = append(wc, "level = :level")
	}

	if filter.PrimaryLanguage != nil {
		data["primary_language"] = *filter.PrimaryLanguage
		wc = append(wc, "primary_language = :primary_language")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE " + strings.Join(wc, " AND "))
	}
}
