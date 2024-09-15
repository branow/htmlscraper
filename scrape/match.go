package scrape

import (
	"strings"
)

type Match func(extract string) (string, bool)

func GetEqualMatch(expected string) Match {
	return func(extract string) (string, bool) {
		return extract, expected == extract
	}
}

func GetPrefixMatch(prefix string) Match {
	return func(extract string) (string, bool) {
		if strings.HasPrefix(extract, prefix) {
			extract = strings.Replace(extract, prefix, "", 1)
			return extract, true
		}
		return extract, false
	}
}
