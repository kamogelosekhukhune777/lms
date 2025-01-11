package coursebus

import "github.com/kamogelosekhukhune777/lms/business/sdk/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByProductID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByProductID       = "product_id"
	OrderByCategory        = "category"
	OrderByPrimaryLanguage = "primary_language"
	OrderByLevel           = "level"
	OrderByPrice           = "price"
)
