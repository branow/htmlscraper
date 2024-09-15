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
