package coursedb

import (
	"fmt"

	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/order"
)

var orderByFields = map[string]string{
	coursebus.OrderByProductID:       "course_id",
	coursebus.OrderByCategory:        "category",
	coursebus.OrderByPrimaryLanguage: "primary_language",
	coursebus.OrderByLevel:           "level",
	coursebus.OrderByPrice:           "price",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
