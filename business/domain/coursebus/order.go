package coursebus

import "github.com/kamogelosekhukhune777/lms/business/sdk/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByProductID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByProductID      = "product_id"
	OrderByTitle          = "title"
	OrderByPrice          = "pricing"
	OrderByPriceLowToHigh = "price_low_to_high"
	OrderByPriceHighToLow = "price_high_to_low"
	OrderByTitleAToZ      = "title_a_to_z"
	OrderByTitleZToA      = "title_z_to_a"
)
