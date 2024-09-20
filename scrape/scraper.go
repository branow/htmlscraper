package scrape

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// The tags that let you to specify where the valuable data is and how to
// get it from the [html.Node].
const (
	SelectorTag  = "select"  // jQuery-like selector to find the node
	ExtractorTag = "extract" // extract operation to get useful data from the node
)

type Mode uint

const (
	Strict Mode = iota
	Tolerant
	Silent
)

// Scraper is a struct that contains a method to scrape data from an
// HTML document ([goquery.Document]).
type Scraper struct {

	//Mode can take three states [Strict], [Tolerant], and [Silent].
	// - [Strict] mode assumes that any error caused during scraping is fatal
	// and stops the following scraping.
	// - [Tolerant] mode assumes that scraping should not be prevented but
	// errors where possible and all errors are returned.
	// - [Silent] mode assumes that scraping should not be stopped by errors
	// where possible and these errors are not returned.
	Mode Mode

	// Extractors is a map that matches custom user extractors to extract tags.
	// Do not use reserved extractor tag names and patterns ([TextExtractTag],
	// [AttrExtractTag], and others), otherwise, the default implementation is executed.
	Extractors map[*Match]Extractor
}

// Scrape scrapes the given doc and writes the useful information into o.
//
// o must be a pointer to a string, slice, or struct, otherwise it causes an error.
// Slices and structs both can contain pointers, strings, slices, and structs but
// the end value must be a string.
//
// selector is a jQuery-like selector that specifies a path to nodes
// (is used in [goquery.Selection.Find]). If selector is empty the doc selection
// (it uses [goquery.Document.Selection]) is considered as default.
//
// extract is a value that specifies how to get useful data from the node.
// extract is required only if o is a pointer to a string or slice, in all
// other cases you can leave it empty.
func (scraper Scraper) Scrape(doc *goquery.Document, o any, selector string, extract string) error {
	err := errors.Join(ValidateNotNil(doc, "doc"), ValidateNotNil(o, "o"))
	if err != nil {
		return ScrapeErr{err}
	}

	ot, ov := reflect.TypeOf(o), reflect.ValueOf(o)
	if ot.Kind() != reflect.Pointer {
		return ScrapeErr{KindErr{Var: "o", KindExp: reflect.Pointer, KindAct: ot.Kind()}}
	}
	ote, ove := ot.Elem(), ov.Elem()

	err = scraper.scrapeObject(doc.Selection, ote, ove, selector, extract)
	if err != nil && scraper.Mode != Silent {
		return ScrapeErr{err}
	}
	return nil
}

func (scraper Scraper) scrapeObject(selection *goquery.Selection, ot reflect.Type, ov reflect.Value, selector string, extract string) error {
	switch ot.Kind() {
	case reflect.String:
		return scraper.scrapeString(selection, ov, selector, extract)
	case reflect.Slice:
		return scraper.scrapeSlice(selection, ot, ov, selector, extract)
	case reflect.Struct:
		return scraper.scrapeStruct(selection, ot, ov, selector)
	case reflect.Pointer:
		return scraper.scrapePointer(selection, ot, ov, selector, extract)
	default:
		kinds := []any{reflect.String, reflect.Slice, reflect.Struct, reflect.Pointer}
		return KindErr{Var: "o", KindExp: kinds, KindAct: ot.Kind()}
	}
}

func (scraper Scraper) scrapeString(selection *goquery.Selection, ov reflect.Value, selector string, extract string) error {
	if len(selector) != 0 {
		selection = selection.Find(selector)
	}

	if selection.Size() == 0 {
		return ScrapingErr{Selector: selector, Cause: NoNodesFoundErr{}}
	}

	node := selection.Nodes[0]
	val, err := scraper.toExtract(node, extract)

	if err != nil {
		return ScrapingErr{Selector: selector, Cause: err}
	}

	ov.SetString(val)
	return nil
}

func (scraper Scraper) scrapeSlice(selection *goquery.Selection, ot reflect.Type, ov reflect.Value, selector string, extract string) error {
	if len(selector) != 0 {
		selection = selection.Find(selector)
	}

	if selection.Size() == 0 {
		return ScrapingErr{Selector: selector, Cause: NoNodesFoundErr{}}
	}

	ote := ot.Elem()
	sv := reflect.MakeSlice(ot, 0, 10)

	errs := []error{}
	selection.EachWithBreak(func(i int, selection *goquery.Selection) bool {
		ve := reflect.New(ote).Elem()
		err := scraper.scrapeObject(selection, ote, ve, "", extract)
		if err != nil {
			s := fmt.Sprintf("%s:n(%d)", selector, i)
			err := ScrapingErr{Selector: s, Cause: err}
			errs = append(errs, err)
		}
		sv = reflect.Append(sv, ve)
		return !(err != nil && scraper.Mode == Strict)
	})

	err := errors.Join(errs...)
	if err == nil || scraper.Mode != Strict {
		ov.Set(sv)
	}

	return err
}

func (scraper Scraper) scrapeStruct(selection *goquery.Selection, ot reflect.Type, ov reflect.Value, selector string) error {
	if len(selector) != 0 {
		selection = selection.Find(selector)
	}

	if selection.Size() == 0 {
		return ScrapingErr{Selector: selector, Cause: NoNodesFoundErr{}}
	}

	errs := []error{}
	selection = selection.First()

	for i := 0; i < ov.NumField(); i++ {
		ft, fv := ot.Field(i), ov.Field(i)
		newSelection, _ := ft.Tag.Lookup(SelectorTag)
		newExtractor, _ := ft.Tag.Lookup(ExtractorTag)
		err := scraper.scrapeObject(selection, ft.Type, fv, newSelection, newExtractor)

		if err != nil {
			err := ScrapingErr{Selector: selector, Cause: err}
			if scraper.Mode == Strict {
				return err
			}
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (scraper Scraper) scrapePointer(selection *goquery.Selection, ot reflect.Type, ov reflect.Value, selector, extract string) error {
	if len(selector) != 0 {
		selection = selection.Find(selector)
	}

	if selection.Size() == 0 {
		return ScrapingErr{Selector: selector, Cause: NoNodesFoundErr{}}
	}

	ote := ot.Elem()
	newValue := reflect.New(ote)
	err := scraper.scrapeObject(selection, ote, newValue.Elem(), "", extract)

	if err != nil {
		err = ScrapingErr{Selector: selector, Cause: err}
	}

	ov.Set(newValue)
	return err
}

func (s Scraper) toExtract(node *html.Node, extract string) (string, error) {
	defaultMap := GetExtractorMap()
	for match, extractor := range defaultMap {
		extract, ok := (*match)(extract)
		if ok {
			return extractor(node, extract)
		}
	}
	customMap := s.Extractors
	for match, extractor := range customMap {
		extract, ok := (*match)(extract)
		if ok {
			return extractor(node, extract)
		}
	}
	return "", ExtractTagErr{ExtractTag: extract}
}
