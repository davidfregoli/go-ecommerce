package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"fregoli.dev/go-ecommerce/components"
	"fregoli.dev/go-ecommerce/db/model"
	. "fregoli.dev/go-ecommerce/db/table"
	"fregoli.dev/go-ecommerce/format"
	"fregoli.dev/go-ecommerce/queries"
	"fregoli.dev/go-ecommerce/types"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/gorilla/sessions"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

const url = "file:./database.sqlite"

var store *sessions.CookieStore
var db *sql.DB

func main() {
	// store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store = sessions.NewCookieStore([]byte("CHANGETHIS"))
	var err error
	db, err = sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	staticHandler := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", staticHandler))

	route("/", homeHandler, middlewareWishlist, middlewareCategories, middlewareSession)
	route("/products/", storeHandler, middlewareWishlist, middlewareCategories, middlewareSession)
	route("/products/{category}/", storeHandler, middlewareWishlist, middlewareCategories, middlewareSession)
	route("/products/{category}/{subcategory}/", subCategoryRedir)
	route("/products/{category}/{subcategory}/{brand}/", brandRedir)
	route("/products/{category}/{subcategory}/{brand}/{product}/{rest...}", productHandler, middlewareWishlist, middlewareCategories, middlewareSession)

	route("POST /wishlist/{action}", wishlistActionHandler, middlewareSession)
	route("POST /search/", searchHandler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func middlewareSession(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := store.Get(r, "session")
		if session.Values["authenticated"] == nil {
			session.Values["authenticated"] = false
			session.Save(r, w)
		}
		ctx = context.WithValue(ctx, types.Authenticated, session.Values["authenticated"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func middlewareCategories(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		stmt := SELECT(
			Category.AllColumns,
		).FROM(
			Category,
		).ORDER_BY(
			CASE().WHEN(Category.Parent.IS_NULL()).THEN(Category.ID).ELSE(Category.Parent).ASC(),
			Category.Parent.ASC(),
			Category.Name.ASC(),
		)
		categories := []model.Category{}
		stmt.Query(db, &categories)
		var parent types.CategoryTree
		categoriesTree := []types.CategoryTree{}
		for i, category := range categories {
			if category.Parent == nil {
				if i != 0 {
					categoriesTree = append(categoriesTree, parent)
				}
				parent = types.CategoryTree{
					Category: category,
					Children: []model.Category{},
				}
			} else {
				parent.Children = append(parent.Children, category)
			}
		}
		ctx = context.WithValue(ctx, types.Categories, categoriesTree)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func middlewareWishlist(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		stmt := SELECT(
			Wishlist.Product,
		).FROM(
			Wishlist,
		).WHERE(
			// Retrieve user from context
			Wishlist.User.EQ(Int(1)),
		)
		wishlist := []int32{}
		stmt.Query(db, &wishlist)
		ctx = context.WithValue(ctx, types.Wishlist, wishlist)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func wishlistActionHandler(w http.ResponseWriter, r *http.Request) {
	action := r.PathValue("action")
	if action != "add" && action != "remove" {
		fmt.Fprintf(w, "wishlistActionHandler() err: %v", "invalid wishlist action")
		return
	}
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	product, err := strconv.Atoi(r.FormValue("product"))
	if err != nil {
		fmt.Fprintf(w, "strconv.Atoi() err: %v", err)
		return
	}

	model := model.Wishlist{
		User:    1,
		Product: int32(product),
	}

	if action == "add" {
		stmt := Wishlist.INSERT(
			Wishlist.User,
			Wishlist.Product,
		).MODEL(model)
		_, err := stmt.Exec(db)
		if err != nil {
			log.Printf("wishlistActionHandler() err: %v\n", err)
		}
	}
	if action == "remove" {
		stmt := Wishlist.DELETE().WHERE(Wishlist.User.EQ(Int(1)).AND(Wishlist.Product.EQ(Int(int64(product)))))
		_, err := stmt.Exec(db)
		if err != nil {
			log.Printf("wishlistActionHandler() err: %v\n", err)
		}
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	search := r.FormValue("query")
	category, err := strconv.ParseInt(r.FormValue("category"), 10, 64)
	if err != nil {
		fmt.Fprintf(w, "strconv.Atoi() err: %v", err)
		return
	}
	stmt := SELECT(Category.AllColumns).FROM(Category).WHERE(Category.ID.EQ(Int(category)))
	var cat model.Category
	stmt.Query(db, &cat)
	if cat.Parent == nil {
		http.Redirect(w, r, "/products/"+cat.Name+"/?search="+search, http.StatusFound)
		return
	}
	stmt = SELECT(Category.AllColumns).FROM(Category).WHERE(Category.ID.EQ(Int(int64(*cat.Parent))))
	var parent model.Category
	stmt.Query(db, &parent)
	http.Redirect(w, r, "/products/"+parent.Name+"/?subCategory="+format.S(cat.ID)+"&search="+search, http.StatusFound)
}

func subCategoryRedir(w http.ResponseWriter, r *http.Request) {
	categorySlug := r.PathValue("category")
	subCategorySlug := r.PathValue("subcategory")
	var subCategory model.Category
	stmt := SELECT(Category.AllColumns).FROM(Category).WHERE(Category.Name.EQ(String(subCategorySlug)))
	stmt.Query(db, &subCategory)
	if subCategory.Parent == nil {
		http.Redirect(w, r, "/products/", http.StatusFound)
		return
	}
	stmt = SELECT(Category.AllColumns).FROM(Category).WHERE(Category.Name.EQ(String(categorySlug)).AND(Category.ID.EQ(Int(int64(*subCategory.Parent)))))
	var category model.Category
	stmt.Query(db, &category)
	http.Redirect(w, r, "/products/"+category.Name+"/?subCategory="+format.S(subCategory.ID), http.StatusFound)
}

func brandRedir(w http.ResponseWriter, r *http.Request) {
	categorySlug := r.PathValue("category")
	subCategorySlug := r.PathValue("subcategory")
	brandSlug := r.PathValue("brand")
	var subCategory model.Category
	stmt := SELECT(Category.AllColumns).FROM(Category).WHERE(Category.Name.EQ(String(subCategorySlug)))
	stmt.Query(db, &subCategory)
	if subCategory.Parent == nil {
		http.Redirect(w, r, "/products/", http.StatusFound)
		return
	}
	stmt = SELECT(Category.AllColumns).FROM(Category).WHERE(Category.Name.EQ(String(categorySlug)).AND(Category.ID.EQ(Int(int64(*subCategory.Parent)))))
	var category model.Category
	stmt.Query(db, &category)
	stmt = SELECT(Brand.AllColumns).FROM(Brand).WHERE(Brand.Slug.EQ(String(brandSlug)))
	var brand model.Brand
	stmt.Query(db, &brand)
	http.Redirect(w, r, "/products/"+category.Name+"/?subCategory="+format.S(subCategory.ID)+"&brand="+format.S(brand.ID), http.StatusFound)
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	data := types.NewData(r.Context())
	data.Page = "store"
	params := r.URL.Query()
	category := r.PathValue("category")
	page := params["page"]
	search := params["search"]
	if len(page) == 0 {
		page = []string{"1"}
	}
	WhereClause := Bool(true)
	ParentCategory := Category.AS("ParentCategory")
	if category == "" {
		data.Store.Categories = queries.SelectRootCategories(db)
	} else {
		var exists bool
		data.Store.Categories, exists = queries.SelectSubCategories(category, db)
		if !exists {
			http.Redirect(w, r, "/products/", http.StatusMovedPermanently)
			return
		}
		data.Store.Category = category
		WhereClause = WhereClause.AND(ParentCategory.Name.EQ(String(data.Store.Category)))
	}
	data.Store.QueryParentCategory = format.MapStringInt(params["parentCategory"])
	data.Store.QueryCategory = format.MapStringInt(params["subCategory"])
	data.Store.QueryBrand = format.MapStringInt(params["brand"])
	data.Store.OrderBy = format.Param(params["orderBy"], "name")
	data.Store.Page = format.ParamInt(page)
	if len(data.Store.QueryParentCategory) > 0 {
		in := []Expression{}
		for _, id := range data.Store.QueryParentCategory {
			in = append(in, Int(id))
		}
		WhereClause = WhereClause.AND(Category.Parent.IN(in...))
	}
	data.Store.Search = format.Param(search, "")
	if data.Store.Search != "" {
		WhereClause = WhereClause.AND(Product.Name.LIKE(String("%" + data.Store.Search + "%")).OR(Brand.Name.LIKE(String("%" + data.Store.Search + "%"))))
	}
	FromClause := Product.INNER_JOIN(
		Brand, Product.Brand.EQ(Brand.ID),
	).INNER_JOIN(
		Category, Product.Category.EQ(Category.ID),
	).INNER_JOIN(
		ParentCategory, Category.Parent.EQ(ParentCategory.ID),
	)
	stmt := SELECT(Brand.AllColumns).FROM(FromClause).WHERE(WhereClause).GROUP_BY(Brand.ID).ORDER_BY(Brand.Name)
	stmt.Query(db, &data.Store.Brands)
	if len(data.Store.QueryBrand) > 0 {
		in := []Expression{}
		for _, id := range data.Store.QueryBrand {
			in = append(in, Int(id))
		}
		WhereClause = WhereClause.AND(Brand.ID.IN(in...))
	}
	if len(data.Store.QueryCategory) > 0 {
		in := []Expression{}
		for _, id := range data.Store.QueryCategory {
			in = append(in, Int(id))
		}
		WhereClause = WhereClause.AND(Category.ID.IN(in...))
	}
	DiscountedPrice := Product.Price.MUL(Int(100).SUB(Product.Discount)).DIV(Int(100))
	stmt = SELECT(MAX(DiscountedPrice).AS("max")).FROM(FromClause).WHERE(WhereClause)
	var dest struct{ Max string }
	stmt.Query(db, &dest)
	n, err := strconv.ParseFloat(dest.Max, 64)
	if err != nil {
		log.Printf("Invalid number: %v", dest.Max)
	}
	data.Store.PriceMax = (n / 100) + 1
	data.Store.Price = []float64{format.ParamFloat(params["priceMin"], 0), format.ParamFloat(params["priceMax"], data.Store.PriceMax)}
	PriceExpression := DiscountedPrice.GT_EQ(Int(int64(data.Store.Price[0] * 100))).AND(
		DiscountedPrice.LT_EQ(Int(int64(data.Store.Price[1] * 100))))
	OrderByClause := Product.Name.ASC()
	switch data.Store.OrderBy {
	case "discount":
		OrderByClause = Product.Discount.DESC()
	case "price":
		OrderByClause = DiscountedPrice.ASC()
	}
	stmt = SELECT(
		Product.AllColumns,
		Brand.AllColumns,
		Category.AllColumns,
		ParentCategory.Name.AS("product_full.parent_category_name"),
		COUNT(STAR).OVER().AS("product_full_count.count"),
	).FROM(
		FromClause,
	).WHERE(
		WhereClause.AND(PriceExpression),
	).ORDER_BY(
		OrderByClause,
	).LIMIT(9).OFFSET(int64(9 * (data.Store.Page - 1)))
	var products []types.ProductFullCount
	data.Store.Products = []types.ProductFull{}
	data.Store.URL = strings.Join([]string{r.URL.Path, r.URL.RawQuery}, "?")
	stmt.Query(db, &products)
	for i, product := range products {
		data.Store.Products = append(data.Store.Products, product.Product)
		if i == 0 {
			data.Store.Count = product.Count
		}
	}
	if data.Store.Price[0] > data.Store.PriceMax {
		URL := format.RemovePriceRange(data.Store.URL)
		http.Redirect(w, r, URL, http.StatusFound)
	}
	components.Page(data).Render(r.Context(), w)
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("product")
	product := queries.SelectProduct(slug, db)
	related := queries.SelectRelatedGallery(product, db)
	data := types.NewData(r.Context())
	data.Lists["related"] = related
	data.Page = "product"
	data.Product = product
	components.Page(data).Render(r.Context(), w)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := types.NewData(r.Context())
	data.Page = "home"
	data.Lists["discounted"] = queries.SelectDiscountedGallery(db)
	components.Page(data).Render(r.Context(), w)
}

func route(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	wrapped := handler
	for i := 0; i < len(middlewares); i++ {
		wrapped = middlewares[i](wrapped)
	}
	http.HandleFunc(path, wrapped)
}

type Middleware func(http.HandlerFunc) http.HandlerFunc
