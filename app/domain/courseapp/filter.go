package courseapp

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
)

type queryParams struct {
	ID              string
	Page            string
	Rows            string
	OrderBy         string
	Category        string
	Level           string
	PrimaryLanguage string
}

func parseQueryParams(r *http.Request) queryParams {
	values := r.URL.Query()
	return queryParams{
		ID:              values.Get("id"),
		Page:            values.Get("page"),
		Rows:            values.Get("rows"),
		OrderBy:         values.Get("orderBy"),
		Category:        values.Get("category"),
		Level:           values.Get("level"),
		PrimaryLanguage: values.Get("primary_language"),
	}
}

func parseFilter(qp queryParams) (coursebus.QueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter coursebus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		switch err {
		case nil:
			filter.ID = &id
		default:
			fieldErrors.Add("product_id", err)
		}
	}

	if qp.Category != "" {
		filter.Category = &qp.Category
	}

	if qp.Level != "" {
		filter.Level = &qp.Level
	}

	if qp.PrimaryLanguage != "" {
		filter.PrimaryLanguage = &qp.PrimaryLanguage
	}

	if len(fieldErrors) > 0 {
		return coursebus.QueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}
