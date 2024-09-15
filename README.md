# `htmlscraper` - automated HTML scraping with jQuery-like selectors in Go

[![Build Status](https://github.com/branow/htmlscraper/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/branow/htmlscraper/actions/workflows/go.yml) 
[![Go Report Card](https://goreportcard.com/badge/github.com/branow/htmlscraper)](https://goreportcard.com/report/github.com/branow/htmlscraper) 
[![PkgGoDev](https://pkg.go.dev/badge/github.com/branow/htmlscraper)](https://pkg.go.dev/github.com/branow/htmlscraper)


## Table of Contents

* [Installation](#installation)
* [Examples](#examples)
* [Contributing](#contributing)
* [License](#license)

## Installation

To install `httpscraper`, use `go get`:

    go get github.com/branow/httpscraper

This will then make the following package available to you:

    github.com/branow/httpscraper/scrape

To update `httpscraper` to the latest version, use `go get -u github.com/branow/httpscraper`.

We currently support the most recent major Go versions from `1.23` onward.

## Examples

Lets scrape the following body tag of [catalog.html](https://github.com/branow/htmlscraper/blob/main/examples/catalog.html) file.

 - [Scrape the catalog name](#scrape-the-catalog-name)
 - [Scrape the product names](#scrape-the-product-names)
 - [Scrape the products](#scrape-the-products)
 - [Scrape the catalog](#scrape-the-catalog)

```html
<body>
    <div class="container">
        <h1>Product Catalog</h1>
        <div class="catalog">
            <div class="product">
                <img src="https://via.placeholder.com/200" alt="Product 1">
                <h2>Product 1</h2>
                <p>Great product for your needs.</p>
                <p class="price">$29.99</p>
            </div>
            <div class="product">
                <img src="https://via.placeholder.com/200" alt="Product 2">
                <h2>Product 2</h2>
                <p>Top-rated product with excellent reviews.</p>
                <p class="price">$39.99</p>
            </div>
            <div class="product">
                <img src="https://via.placeholder.com/200" alt="Product 3">
                <h2>Product 3</h2>
                <p>Best value for your money.</p>
                <p class="price">$19.99</p>
            </div>
        </div>
    </div>
</body>
```

### Scrape the catalog name

```Go
package examples

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/branow/htmlscraper/scrape"
)

func ScrapeString() {
	// create goquery document
	file := getCatalogFile() //get catalog.html
	defer file.Close()
	doc, err := goquery.NewDocumentFromReader(file)
	raisePanic(err)

	// create Scraper
	scraper := scrape.Scraper{}

	// scraping
	var catalog string //product catalog name
	err = scraper.Scrape(doc, &catalog, ".container > h1", "text")

	// get output
	fmt.Println("Got Error:", err)
	fmt.Println("Got Output:", catalog)
}

```
It prints:
```
Got Error: <nil>
Got Output: Product Catalog
```
[The example file.](https://github.com/branow/htmlscraper/blob/main/examples/scrape_string.go)

### Scrape the product names

```go
package examples

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/branow/htmlscraper/scrape"
)

func ScrapeSliceOfStrings() {
	// create goquery document
	file := getCatalogFile()
	defer file.Close()
	doc, err := goquery.NewDocumentFromReader(file)
	raisePanic(err)

	// create Scraper
	scraper := scrape.Scraper{}

	// scraping
	var products []string //product names
	err = scraper.Scrape(doc, &products, ".product > h2", "text")

	// get output
	fmt.Println("Got Error:", err)
	fmt.Println("Got Output:", products)
}

```
It prints:
```
Got Error: <nil>
Got Output: [Product 1 Product 2 Product 3]
```
[The example file.](https://github.com/branow/htmlscraper/blob/main/examples/scrape_slice_of_strings.go)

### Scrape the products

```go
package examples

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/branow/htmlscraper/scrape"
	"golang.org/x/net/html"
)

func ScrapeSliceOfStructs() {
	// create goquery document
	file := getCatalogFile()
	defer file.Close()
	doc, err := goquery.NewDocumentFromReader(file)
	raisePanic(err)

	// create custom extractor for price data
	priceMatch := scrape.GetEqualMatch("*price")
	priceExtractor := func(node *html.Node, extract string) (string, error) {
		price := node.FirstChild.Data
		return strings.Replace(price, "$", "", 1), nil
	}
	customExtractors := map[*scrape.Match]scrape.Extractor{&priceMatch: priceExtractor}

	// create Scraper
	scraper := scrape.Scraper{Extractors: customExtractors}

	// scraping
	type Product struct {
		Name        string `select:"h2" extract:"text"`
		Description string `select:"p" extract:"text"`
		Price       string `select:".price" extract:"*price"`
		Image       string `select:"img" extract:"@src"`
	}
	var products []Product
	err = scraper.Scrape(doc, &products, ".product", "")

	// get output
	fmt.Println("Got Error:", err)
	fmt.Println("Got Output:")
	for _, p := range products {
		fmt.Println(p)
	}
}

```
It prints:
```
Got Error: <nil>
Got Output:
{Product 1 Great product for your needs. 29.99 https://via.placeholder.com/200}
{Product 2 Top-rated product with excellent reviews. 39.99 https://via.placeholder.com/200}
{Product 3 Best value for your money. 19.99 https://via.placeholder.com/200}
```
[The example file.](https://github.com/branow/htmlscraper/blob/main/examples/scrape_slice_of_structs.go)

### Scrape the catalog

```go
package examples

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/branow/htmlscraper/scrape"
	"golang.org/x/net/html"
)

func ScrapeStruct() {
	// create goquery document
	file := getCatalogFile()
	defer file.Close()
	doc, err := goquery.NewDocumentFromReader(file)
	raisePanic(err)

	// create custom extractor for price data
	priceMatch := scrape.GetEqualMatch("*price")
	priceExtractor := func(node *html.Node, extract string) (string, error) {
		price := node.FirstChild.Data
		return strings.Replace(price, "$", "", 1), nil
	}
	customExtractors := map[*scrape.Match]scrape.Extractor{&priceMatch: priceExtractor}

	// create Scraper
	scraper := scrape.Scraper{Extractors: customExtractors}

	// scraping
	type Product struct {
		Name        string `select:"h2" extract:"text"`
		Description string `select:"p" extract:"text"`
		Price       string `select:".price" extract:"*price"`
		Image       string `select:"img" extract:"@src"`
	}
	type Catalog struct {
		Name     string    `select:"h1" extract:"text"`
		Products []Product `select:".product"`
	}
	var catalog Catalog
	err = scraper.Scrape(doc, &catalog, ".container", "")

	// get output
	fmt.Println("Got Error:", err)
	fmt.Println("Got Output:")
	fmt.Println("Catalog {")
	fmt.Println(catalog.Name)
	for _, p := range catalog.Products {
		fmt.Println(p)
	}
	fmt.Println("}")
}

```
It prints:
```
Catalog {
Product Catalog
{Product 1 Great product for your needs. 29.99 https://via.placeholder.com/200}
{Product 2 Top-rated product with excellent reviews. 39.99 https://via.placeholder.com/200}
{Product 3 Best value for your money. 19.99 https://via.placeholder.com/200}
}
```
[The example file.](https://github.com/branow/htmlscraper/blob/main/examples/scrape_struct.go)

## Contributing

Please feel free to submit issues, fork the repository and send pull requests!

## License

This project is licensed under the terms of the [MIT license](https://github.com/branow/htmlscraper/blob/main/LICENSE.txt).
