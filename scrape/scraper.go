package scrape

import (
	"errors"
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

// Scraper is a struct that contains a method to scrape data from an
// HTML document ([goquery.Document]).
type Scraper struct {

	// If Strict flag is true, the [Scraper.Scrape] method
	// returns an error if the seeking HTML node is not found otherwise
	// it returns a zero value according to the type. The exception is an
	// slice type for which the flag does not work and even if there
	// are not found notes it returns an empty slice of the specified type
	// with a capacity of 10.
	Strict bool

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
	err := errors.Join(checkNil(doc, "doc"), checkNil(o, "o"))
	if err != nil {
		return err
	}

	ot, ov := reflect.TypeOf(o), reflect.ValueOf(o)
	if ot.Kind() != reflect.Pointer {
		return GetKindErr(ot, reflect.Pointer, ot.Kind())
	}
	ote, ove := ot.Elem(), ov.Elem()
	return scraper.scrapeObject(doc.Selection, ote, ove, selector, extract)
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
		return GetMultiKindErr(ot, []any{reflect.String, reflect.Slice, reflect.Struct}, ot.Kind())
	}
}

func (scraper Scraper) scrapeString(selection *goquery.Selection, ov reflect.Value, selector string, extract string) error {
	if len(selector) != 0 {
		selection = selection.Find(selector)
	}
	if selection.Size() == 0 {
		if scraper.Strict {
			return GetNotFoundErr(selector)
		} else {
			ov.SetString("")
			return nil
		}
	}

	node := selection.Nodes[0]
	val, err := scraper.toExtract(node, extract)
	if err != nil {
		if err.Error() == GetExtractErr(extract).Error() {
			return err
		}
		return WrapExtractErr(selector, err)
	}
	ov.SetString(val)
	return nil
}

func (scraper Scraper) scrapeSlice(selection *goquery.Selection, ot reflect.Type, ov reflect.Value, selector string, extract string) error {
	ote := ot.Elem()
	sv := reflect.MakeSlice(ot, 0, 10)
	errs := []error{}
	selection.Find(selector).Each(func(i int, selection *goquery.Selection) {
		ve := reflect.New(ote).Elem()
		err := scraper.scrapeObject(selection, ote, ve, "", extract)
		if err != nil {
			errs = append(errs, err)
		}
		sv = reflect.Append(sv, ve)
	})
	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	ov.Set(sv)
	return nil
}

func (scraper Scraper) scrapeStruct(selection *goquery.Selection, ot reflect.Type, ov reflect.Value, selector string) error {
	if len(selector) != 0 {
		selection = selection.Find(selector)
	}
	if selection.Size() == 0 {
		if scraper.Strict {
			return GetNotFoundErr(selector)
		} else {
			return nil
		}
	}
	selection = selection.First()
	for i := 0; i < ov.NumField(); i++ {
		ft, fv := ot.Field(i), ov.Field(i)
		selector, _ := ft.Tag.Lookup(SelectorTag)
		extractor, _ := ft.Tag.Lookup(ExtractorTag)
		err := scraper.scrapeObject(selection, ft.Type, fv, selector, extractor)
		if err != nil {
			return err
		}
	}
	return nil
}

func (scraper Scraper) scrapePointer(selection *goquery.Selection, ot reflect.Type, ov reflect.Value, selector, extract string) error {
	if len(selector) != 0 {
		selection = selection.Find(selector)
	}
	if selection.Size() == 0 {
		if scraper.Strict {
			return GetNotFoundErr(selector)
		} else {
			return nil
		}
	}
	ote := ot.Elem()
	newValue := reflect.New(ote)
	err := scraper.scrapeObject(selection, ote, newValue.Elem(), "", extract)
	if err != nil {
		return err
	}
	ov.Set(newValue)
	return nil
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
	return "", GetExtractErr(extract)
}
