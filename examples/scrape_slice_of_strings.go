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
