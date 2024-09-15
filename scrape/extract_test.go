package scrape_test

import (
	"bytes"
	"testing"

	"github.com/PuerkitoBio/goquery"
	. "github.com/branow/htmlscraper/scrape"
	"github.com/branow/tabtest/tab"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestDeepExtractText(t *testing.T) {
	htmldata := `<span class="x">The <span class="x"><span class="y">government</span></span> has <span class="cl">set a</span> growth <span class="cl">target</span> of 6%.</span>`
	doc, _ := goquery.NewDocumentFromReader(bytes.NewBufferString(htmldata))
	node := doc.Find(".x").First().Nodes[0]
	exp := "The government has set a growth target of 6%."
	act := ExtractDeepText(node)
	assert.Equal(t, exp, act)
}

func TestExtractText(t *testing.T) {
	htmldata := `<span class="x">The <span class="x"><span class="y">government</span></span> has <span class="cl">set a</span> growth <span class="cl">target</span> of 6%.</span>`
	doc, _ := goquery.NewDocumentFromReader(bytes.NewBufferString(htmldata))
	node := doc.Find(".x").First().Nodes[0]
	exp := "The  has  growth  of 6%."
	act := ExtractText(node)
	assert.Equal(t, exp, act)
}

func TestExtractAttribute(t *testing.T) {
	args := []tab.Args{
		{
			"@attr not found",
			&html.Node{
				Attr: []html.Attribute{{Namespace: "", Key: "src", Val: "www.site.com"}},
			},
			"source",
			"",
			GetAttributeNotFoundErr("source"),
		},
		{
			"@attr source",
			&html.Node{
				Attr: []html.Attribute{{Namespace: "", Key: "source", Val: "www.site.com"}},
			},
			"source",
			"www.site.com",
			nil,
		},
	}
	test := func(t *testing.T, node *html.Node, attr, exp string, eErr error) {
		act, aErr := ExtractAttribute(node, attr)

		if eErr == nil {
			assert.NoError(t, aErr)
		} else {
			if assert.Error(t, aErr) {
				assert.EqualError(t, aErr, eErr.Error())
			}
		}

		assert.Equal(t, exp, act)
	}
	tab.RunWithArgs(t, args, test)
}
