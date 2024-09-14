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
	priceExtractor := func(node *html.Node) (string, error) {
		price := node.FirstChild.Data
		return strings.Replace(price, "$", "", 1), nil
	}
	customExtractors := map[string]scrape.Extractor{"*price": priceExtractor}

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
