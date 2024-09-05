package format

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/a-h/templ"
)

var S = fmt.Sprint
var F = fmt.Sprintf

func USD(amt int32) string {
	f := float64(amt) / 100
	return fmt.Sprintf("$%.2f", f)
}

func USDD(amt int32, dsct int32) string {
	inverse := 100 - dsct
	price := amt * inverse
	f := float64(price) / 10000
	return fmt.Sprintf("$%.2f", f)
}

func MapStringInt(data []string) []int64 {
	mapped := make([]int64, len(data))
	for i, el := range data {
		n, err := strconv.ParseInt(el, 10, 64)
		if err != nil {
			log.Printf("Invalid number: %v", el)
			continue
		}
		mapped[i] = n
	}
	return mapped
}

func ParamFloat(data []string, fallback float64) float64 {
	if len(data) != 1 {
		return fallback
	}
	val := data[0]
	n, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Printf("Invalid number: %v", val)
		return fallback
	}
	return n
}

func ParamInt(data []string) int {
	if len(data) != 1 {
		log.Printf("Multiple page values: %v", data)
		return 1
	}
	val := data[0]
	n, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Invalid number: %v", val)
		return 1
	}
	return n
}

func Param(data []string, fallback string) string {
	if len(data) != 1 {
		return fallback
	}
	return data[0]
}

func URL(parts ...string) templ.SafeURL {
	url := "/"
	for _, part := range parts {
		url += part
		url += "/"
	}
	return templ.URL(url)
}

func QURL(url templ.SafeURL, query string) templ.SafeURL {
	return url + templ.SafeURL(query)
}

func PageURL(current string, page int) templ.SafeURL {
	url, err := url.Parse(current)
	if err != nil {
		log.Fatal(err)
	}
	values := url.Query()
	values.Set("page", S(page))
	url.RawQuery = values.Encode()
	unsafe := url.String()
	return templ.URL(unsafe)
}

func RemovePriceRange(current string) string {
	url, err := url.Parse(current)
	if err != nil {
		log.Fatal(err)
	}
	values := url.Query()
	values.Del("priceMin")
	values.Del("priceMax")
	url.RawQuery = values.Encode()
	return url.String()
}

func URLParams(current string) map[string][]string {
	url, err := url.Parse(current)
	if err != nil {
		log.Fatal(err)
	}
	values := url.Query()
	return values
}

func CategoryURL(current string, category string) templ.SafeURL {
	url, err := url.Parse(current)
	if err != nil {
		log.Fatal(err)
	}
	values := url.Query()
	url.Path = "/products/" + category + "/"
	values.Del("parentCategory")
	values.Del("page")
	url.RawQuery = values.Encode()
	unsafe := url.String()
	return templ.URL(unsafe)
}
