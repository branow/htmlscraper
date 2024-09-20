package scrape

import (
	"strings"

	"golang.org/x/net/html"
)

// Extractor tags to specify extract operations.
const (
	TextExtractTag     = "text"     // get a text of children's text nodes
	DeepTextExtractTag = "deeptext" // get a text of descendants' text nodes
	AttrExtractTag     = "@"        // get a value of an attribute ("@href", "@src")
)

// Extractor is a function that processes the given node and returns
// the valuable data in string format.
type Extractor func(node *html.Node, extract string) (string, error)

// GetExtractorMap returns the default map to match extracting tags
// and extracting functions (or extractors).
func GetExtractorMap() map[*Match]Extractor {
	m := map[*Match]Extractor{}

	textMatch := GetEqualMatch(TextExtractTag)
	m[&textMatch] = func(node *html.Node, extract string) (string, error) {
		data := ExtractText(node)
		return data, nil
	}

	deepTextMatch := GetEqualMatch(DeepTextExtractTag)
	m[&deepTextMatch] = func(node *html.Node, extract string) (string, error) {
		data := ExtractDeepText(node)
		return data, nil
	}

	attrMatch := GetPrefixMatch(AttrExtractTag)
	m[&attrMatch] = ExtractAttribute

	return m
}

// ExtractDeepText returns the text of all descendants' text nodes.
func ExtractDeepText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	n := node.FirstChild
	text := []string{}
	for n != nil {
		t := ExtractDeepText(n)
		text = append(text, t)
		n = n.NextSibling
	}
	return strings.Join(text, "")
}

// ExtractDeepText returns the text of all children's text nodes.
func ExtractText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	n := node.FirstChild
	text := []string{}
	for n != nil {
		if n.Type == html.TextNode {
			text = append(text, n.Data)
		}
		n = n.NextSibling
	}
	return strings.Join(text, "")
}

// ExtractAttribute returns the value of the given attribute.
// If the attribute is absent it returns an error.
func ExtractAttribute(node *html.Node, attr string) (string, error) {
	for _, v := range node.Attr {
		if v.Key == attr {
			return v.Val, nil
		}
	}
	return "", AttributeNotFoundErr{Attr: attr}
}
