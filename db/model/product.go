//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type Product struct {
	ID          int32 `sql:"primary_key"`
	Slug        string
	Name        string
	Description string
	Category    int32
	Asin        string
	Price       int32
	Discount    int32
	Brand       int32
}
