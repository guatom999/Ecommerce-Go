package products

import (
	"github.com/guatom999/Ecommerce-Go/modules/appinfo"
	"github.com/guatom999/Ecommerce-Go/modules/entities"
)

type Product struct {
	Id          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Category    *appinfo.Category `json:"category"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	Price       float64           `json:"price"`
	Image       []*entities.Image `json:"images"`
}

type ProductFilter struct {
	Id     string `query:"id"`
	Search string `query:"search"`
	*entities.PaginationReq
	*entities.SortReq
}
