//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/sqlite"
)

var Product = newProductTable("", "product", "")

type productTable struct {
	sqlite.Table

	// Columns
	ID          sqlite.ColumnInteger
	Slug        sqlite.ColumnString
	Name        sqlite.ColumnString
	Description sqlite.ColumnString
	Category    sqlite.ColumnInteger
	Asin        sqlite.ColumnString
	Price       sqlite.ColumnInteger
	Discount    sqlite.ColumnInteger
	Brand       sqlite.ColumnInteger

	AllColumns     sqlite.ColumnList
	MutableColumns sqlite.ColumnList
}

type ProductTable struct {
	productTable

	EXCLUDED productTable
}

// AS creates new ProductTable with assigned alias
func (a ProductTable) AS(alias string) *ProductTable {
	return newProductTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ProductTable with assigned schema name
func (a ProductTable) FromSchema(schemaName string) *ProductTable {
	return newProductTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ProductTable with assigned table prefix
func (a ProductTable) WithPrefix(prefix string) *ProductTable {
	return newProductTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ProductTable with assigned table suffix
func (a ProductTable) WithSuffix(suffix string) *ProductTable {
	return newProductTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newProductTable(schemaName, tableName, alias string) *ProductTable {
	return &ProductTable{
		productTable: newProductTableImpl(schemaName, tableName, alias),
		EXCLUDED:     newProductTableImpl("", "excluded", ""),
	}
}

func newProductTableImpl(schemaName, tableName, alias string) productTable {
	var (
		IDColumn          = sqlite.IntegerColumn("id")
		SlugColumn        = sqlite.StringColumn("slug")
		NameColumn        = sqlite.StringColumn("name")
		DescriptionColumn = sqlite.StringColumn("description")
		CategoryColumn    = sqlite.IntegerColumn("category")
		AsinColumn        = sqlite.StringColumn("asin")
		PriceColumn       = sqlite.IntegerColumn("price")
		DiscountColumn    = sqlite.IntegerColumn("discount")
		BrandColumn       = sqlite.IntegerColumn("brand")
		allColumns        = sqlite.ColumnList{IDColumn, SlugColumn, NameColumn, DescriptionColumn, CategoryColumn, AsinColumn, PriceColumn, DiscountColumn, BrandColumn}
		mutableColumns    = sqlite.ColumnList{SlugColumn, NameColumn, DescriptionColumn, CategoryColumn, AsinColumn, PriceColumn, DiscountColumn, BrandColumn}
	)

	return productTable{
		Table: sqlite.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:          IDColumn,
		Slug:        SlugColumn,
		Name:        NameColumn,
		Description: DescriptionColumn,
		Category:    CategoryColumn,
		Asin:        AsinColumn,
		Price:       PriceColumn,
		Discount:    DiscountColumn,
		Brand:       BrandColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}