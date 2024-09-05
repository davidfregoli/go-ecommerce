package queries

import (
	"database/sql"

	"fregoli.dev/go-ecommerce/db/model"
	. "fregoli.dev/go-ecommerce/db/table"
	"fregoli.dev/go-ecommerce/types"
	. "github.com/go-jet/jet/v2/sqlite"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func SelectProduct(slug string, db *sql.DB) types.ProductFull {
	var product types.ProductFull
	ParentCategory := Category.AS("ParentCategory")
	stmt := SELECT(
		Product.AllColumns,
		Brand.AllColumns,
		Category.AllColumns,
		ParentCategory.Name.AS("product_full.parent_category_name"),
	).FROM(
		Product.INNER_JOIN(
			Brand, Product.Brand.EQ(Brand.ID),
		).INNER_JOIN(
			Category, Product.Category.EQ(Category.ID),
		).INNER_JOIN(
			ParentCategory, Category.Parent.EQ(ParentCategory.ID),
		),
	).WHERE(
		Product.Slug.EQ(String(slug)),
	)
	stmt.Query(db, &product)
	return product
}

func SelectDiscountedGallery(db *sql.DB) []types.ProductFull {
	ParentCategory := Category.AS("ParentCategory")
	stmt := SELECT(
		Product.AllColumns,
		Brand.AllColumns,
		Category.AllColumns,
		ParentCategory.Name.AS("product_full.parent_category_name"),
	).FROM(
		Product.INNER_JOIN(
			Brand, Product.Brand.EQ(Brand.ID),
		).INNER_JOIN(
			Category, Product.Category.EQ(Category.ID),
		).INNER_JOIN(
			ParentCategory, Category.Parent.EQ(ParentCategory.ID),
		),
	).WHERE(
		Product.Discount.GT(Int(0)),
	).ORDER_BY(
		Product.Discount.DESC(),
	).LIMIT(8)
	var discounted []types.ProductFull
	stmt.Query(db, &discounted)
	return discounted
}

func SelectRelatedGallery(product types.ProductFull, db *sql.DB) []types.ProductFull {
	ParentCategory := Category.AS("ParentCategory")
	stmt := SELECT(
		Product.AllColumns,
		Brand.AllColumns,
		Category.AllColumns,
		ParentCategory.Name.AS("product_full.parent_category_name"),
	).FROM(
		Product.INNER_JOIN(
			Brand, Product.Brand.EQ(Brand.ID),
		).INNER_JOIN(
			Category, Product.Category.EQ(Category.ID),
		).INNER_JOIN(
			ParentCategory, Category.Parent.EQ(ParentCategory.ID),
		),
	).WHERE(
		Product.Category.EQ(Int(int64(product.Category.ID))).AND(Product.ID.NOT_EQ(Int(int64(product.ID)))),
	).ORDER_BY(
		Raw("RANDOM()"),
	).LIMIT(4)
	var discounted []types.ProductFull
	stmt.Query(db, &discounted)
	return discounted
}

func SelectRootCategories(db *sql.DB) []model.Category {
	var categories []model.Category
	stmt := SELECT(Category.AllColumns).FROM(Category).WHERE(Category.Parent.IS_NULL()).ORDER_BY(Category.Name.ASC())
	stmt.Query(db, &categories)
	return categories
}

func SelectSubCategories(parent string, db *sql.DB) ([]model.Category, bool) {
	ParentCategory := Category.AS("ParentCategory")
	var categories []model.Category
	stmt := SELECT(Category.AllColumns).FROM(Category.INNER_JOIN(
		ParentCategory, Category.Parent.EQ(ParentCategory.ID),
	)).WHERE(ParentCategory.Name.EQ(String(parent))).ORDER_BY(Category.Name.ASC())
	stmt.Query(db, &categories)
	return categories, len(categories) > 0
}

func SelectStoreProducts(db *sql.DB) []model.Product {
	var products []model.Product
	stmt := SELECT(Product.AllColumns, COUNT(STAR).OVER().AS("count")).FROM(Product).LIMIT(9)
	stmt.Query(db, &products)
	// fmt.Println(stmt.DebugSql())
	return products
}
