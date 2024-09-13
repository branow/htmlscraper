package scrape

import (
	"errors"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const (
	SelectorTag  = "select"
	ExtractorTag = "extract"
)

const (
	TextExtractTag = "text"
	AttrExtractTag = "@"
	FuncExtractTag = "()"
)

type Extractor func(node *html.Node) (string, error)

type Scraper struct {
	Strict     bool
	Extractors map[string]Extractor
}

func (scraper Scraper) Scrape(doc *goquery.Document, o any, selector string, extract string) error {
	err := errors.Join(checkNil(doc, "doc"), checkNil(o, "variable"))
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

	extractor, err := scraper.getExtractor(extract)
	if err != nil {
		return err
	}

	node := selection.Nodes[0]
	val, err := extractor(node)
	if err != nil {
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

func (s Scraper) getExtractor(extract string) (Extractor, error) {
	if extract == TextExtractTag {
		return func(node *html.Node) (string, error) {
			return node.FirstChild.Data, nil
		}, nil
	} else if strings.HasPrefix(extract, AttrExtractTag) {
		attr := extract[1:]
		return func(node *html.Node) (string, error) {
			for _, v := range node.Attr {
				if v.Key == attr {
					return v.Val, nil
				}
			}
			return "", GetAttributeNotFoundErr(attr)
		}, nil
	} else if extractor, ok := s.Extractors[extract]; ok {
		return extractor, nil
	}
	return nil, GetExtractErr(extract)
}
