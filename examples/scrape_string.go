package examples

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/branow/htmlscraper/scrape"
)

func ScrapeString() {
	// create goquery document
	file := getCatalogFile()
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
