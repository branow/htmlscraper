package scrape

import (
	"strings"
)

// Match wraps boolean logic of matching values of
// extracting tags with extracting function (or extractors).
// Match returns already processed value of extracting tag.
// (an example "@href" -> "href").
type Match func(extract string) (string, bool)

// GetEqualMatch creates a Match function that compares
// the given value with the value of the extracting tag.
func GetEqualMatch(expected string) Match {
	return func(extract string) (string, bool) {
		return extract, expected == extract
	}
}

// GetPrefixMatch creates a Match function that checks whether
// the extracting tag value has the given prefix and returns
// a boolean result with the extracting tag value. In true case,
// it cuts the matched prefix from the extracted value
// (an example "@href" -> "href")
func GetPrefixMatch(prefix string) Match {
	return func(extract string) (string, bool) {
		if strings.HasPrefix(extract, prefix) {
			extract = strings.Replace(extract, prefix, "", 1)
			return extract, true
		}
		return extract, false
	}
}
