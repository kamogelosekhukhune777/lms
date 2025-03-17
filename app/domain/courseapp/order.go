package courseapp

import "github.com/kamogelosekhukhune777/lms/business/domain/coursebus"

var orderByFields = map[string]string{
	"product_id":        coursebus.OrderByProductID,
	"title":             coursebus.OrderByTitle,
	"pricing":           coursebus.OrderByPrice,
	"price_low_to_high": coursebus.OrderByPriceLowToHigh,
	"price_high_to_low": coursebus.OrderByPriceHighToLow,
	"title_a_to_z":      coursebus.OrderByTitleAToZ,
	"title_z_to_a":      coursebus.OrderByTitleZToA,
}
