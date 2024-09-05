package types

import (
	"context"

	"fregoli.dev/go-ecommerce/db/model"
)

type Data struct {
	Authenticated bool
	Categories    []CategoryTree
	Lists         map[string][]ProductFull
	Page          string
	Product       ProductFull
	Wishlist      []int32
	Store         Store
}

type CategoryTree struct {
	model.Category
	Children []model.Category
}

func NewData(ctx context.Context) Data {
	return Data{
		Lists:         map[string][]ProductFull{},
		Authenticated: ctx.Value(Authenticated).(bool),
		Categories:    ctx.Value(Categories).([]CategoryTree),
		Wishlist:      ctx.Value(Wishlist).([]int32),
		Store:         Store{},
	}
}

type ProductFull struct {
	model.Product
	Brand              model.Brand
	Category           model.Category
	ParentCategoryName string
}

type ProductFullCount struct {
	Product ProductFull
	Count   int
}

type Store struct {
	Brands              []model.Brand
	Categories          []model.Category
	Category            string
	Count               int
	OrderBy             string
	Page                int
	Price               []float64
	PriceMax            float64
	Products            []ProductFull
	QueryBrand          []int64
	QueryCategory       []int64
	QueryParentCategory []int64
	Search              string
	URL                 string
}

func (s Store) CountPages() int {
	rem := 0
	mod := s.Count % 9
	if mod > 0 {
		rem = 1
	}
	div := s.Count / 9
	return div + rem
}

type ContextKey string

var Authenticated = ContextKey("authenticated")
var Categories = ContextKey("categories")
var Wishlist = ContextKey("wishlist")
