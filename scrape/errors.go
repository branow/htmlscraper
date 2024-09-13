package scrape

import (
	"fmt"
	"strings"
)

func WrapExtractErr(selector string, err error) error {
	err = fmt.Errorf("%w by selector: %s", err, selector)
	return WrapScrapeErr(err)
}

func GetAttributeNotFoundErr(attr string) error {
	return fmt.Errorf("attribute %s not found", attr)
}

func GetNotFoundErr(selector string) error {
	err := fmt.Errorf("no nodes found by selector %s", selector)
	return WrapScrapeErr(err)
}

func GetExtractErr(extract string) error {
	err := fmt.Errorf("invalid extract: %s", extract)
	return WrapScrapeErr(err)
}

func GetMultiKindErr(typeName any, expKinds []any, actKind any) error {
	kinds := mapSlice(expKinds, func(k any) string { return fmt.Sprintf("%v", k) })
	err := fmt.Errorf("%v must be a %s, but it is a %v", typeName, strings.Join(kinds, ","), actKind)
	return WrapScrapeErr(err)
}

func GetKindErr(typeName, expKind, actKind any) error {
	err := fmt.Errorf("%v must be a %v, but it is a %v", typeName, expKind, actKind)
	return WrapScrapeErr(err)
}

func GetNilErr(name string) error {
	return WrapScrapeErr(fmt.Errorf("is nil"))
}

func WrapScrapeErr(err error) error {
	return fmt.Errorf("scrape: %w", err)
}
