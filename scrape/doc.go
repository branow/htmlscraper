// Package scrape implements scraping functionality for extracting
// useful data from HTTP text using jQuery-like selectors.
// It contains struct [Scraper] with method [Scraper.Scrape] built
// on the [github.com/PuerkitoBio.goquery] library.
//
// Here is a simple example, scraping a Product struct from a html
// document.
//
//	 htmlData := `<div class="product">
//	 	<img src="https://via.placeholder.com/200" alt="Product 1">
//	 	<h2>Product 1</h2>
//	 	<p>Great product for your needs.</p>
//	 	<p class="price">$29.99</p>
//	 </div>`
//	 r := bytes.NewBufferString(htmlData)
//	 doc, _ := goquery.NewDocumentFromReader(r)
//
//	 scraper := scrape.Scraper{}
//
//	 // scraping
//	 type Product struct {
//		Name        string `select:"h2" extract:"text"`
//	 	Description string `select:"p" extract:"text"`
//	 	Price       string `select:".price" extract:"text"`
//	 	Image       string `select:"img" extract:"@src"`
//	 }
//	 var product Product
//	 err := scraper.Scrape(doc, &product, ".product", "")
//
//	 // get output
//	 fmt.Println("Got Error:", err)
//	 fmt.Println("Got Output:")
//	 fmt.Println(product)
package scrape
