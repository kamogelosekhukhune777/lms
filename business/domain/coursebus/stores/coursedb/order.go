package coursedb

import (
	"fmt"

	"github.com/kamogelosekhukhune777/lms/business/domain/coursebus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/order"
)

var orderByFields = map[string]string{
	coursebus.OrderByProductID:      "course_id",
	coursebus.OrderByPriceLowToHigh: "pricing",
	coursebus.OrderByPriceHighToLow: "pricing",
	coursebus.OrderByTitleAToZ:      "title",
	coursebus.OrderByTitleZToA:      "title",
}

func orderByClause(orderBy order.By) (string, error) {
	// Validate field existence
	if _, exists := orderByFields[orderBy.Field]; !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	// Ensure order direction is valid
	validDirections := map[string]bool{"ASC": true, "DESC": true}
	if _, valid := validDirections[orderBy.Direction]; !valid {
		return "", fmt.Errorf("invalid order direction %q", orderBy.Direction)
	}

	// Handle predefined order cases
	switch orderBy.Field {
	case coursebus.OrderByPriceLowToHigh:
		return " ORDER BY pricing ASC", nil
	case coursebus.OrderByPriceHighToLow:
		return " ORDER BY pricing DESC", nil
	case coursebus.OrderByTitleAToZ:
		return " ORDER BY title ASC", nil
	case coursebus.OrderByTitleZToA:
		return " ORDER BY title DESC", nil
	default:
		return fmt.Sprintf(" ORDER BY %s %s", orderByFields[orderBy.Field], orderBy.Direction), nil
	}
}
