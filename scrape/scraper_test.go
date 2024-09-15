package scrape_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	. "github.com/branow/htmlscraper/scrape"
	"github.com/branow/tabtest/tab"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

type ScrapeCfg struct {
	CaseName   string
	extractors map[string]Extractor
	strict     bool
	doc        *goquery.Document
	o          any
	selector   string
	extract    string
	exp        any
	eErr       error
}

func test(t *testing.T, c ScrapeCfg) {
	scraper := Scraper{Strict: c.strict}
	if c.extractors != nil {
		scraper.Extractors = c.extractors
	}

	aErr := scraper.Scrape(c.doc, c.o, c.selector, c.extract)

	if c.eErr == nil {
		assert.NoError(t, aErr)
	} else {
		if assert.Error(t, aErr) {
			assert.EqualError(t, aErr, c.eErr.Error())
		}
	}

	assert.Equal(t, c.exp, c.o)
}

func TestScraper_Scrape(t *testing.T) {
	htmldata := `<div class="webtop"><h1 class="headword" id="direct_h_1" opal_spoken="y" random="y" htag="h1" hclass="headword" ox3000="y" opal_written="y">direct</h1> <span class="pos" hclass="pos" htag="span">adjective</span><div class="symbols" hclass="symbols" htag="div"><a href="https://www.oxfordlearnersdictionaries.com/wordlists/oxford3000-5000?dataset=english&amp;list=ox3000&amp;level=a2"><span class="ox3ksym_a2">&nbsp;</span></a><a href="https://www.oxfordlearnersdictionaries.com/wordlists/opal?dataset=english&amp;list=opal_written&amp;level=sublist_3"><span class="opal_symbol" href="OPAL_Written::Sublist_3">OPAL W</span></a><a href="https://www.oxfordlearnersdictionaries.com/wordlists/opal?dataset=english&amp;list=opal_spoken&amp;level=sublist_4"><span class="opal_symbol" href="OPAL_Spoken::Sublist_4">OPAL S</span></a></div><span class="phonetics"> <div class="phons_br" htag="div" geo="br" wd="direct" hclass="phons_br"><div class="sound audio_play_button pron-uk icon-audio" data-src-mp3="https://www.oxfordlearnersdictionaries.com/media/english/uk_pron/d/dir/direc/direct__gb_7.mp3" data-src-ogg="https://www.oxfordlearnersdictionaries.com/media/english/uk_pron_ogg/d/dir/direc/direct__gb_7.ogg" title="direct pronunciation
                    English" style="cursor: pointer" valign="top">&nbsp;</div><span class="phon">/dəˈrekt/</span><span class="sep">,</span> <div class="sound audio_play_button pron-uk icon-audio" data-src-mp3="https://www.oxfordlearnersdictionaries.com/media/english/uk_pron/d/dir/direc/direct__gb_8.mp3" data-src-ogg="https://www.oxfordlearnersdictionaries.com/media/english/uk_pron_ogg/d/dir/direc/direct__gb_8.ogg" title="direct pronunciation
                    English" style="cursor: pointer" valign="top">&nbsp;</div><span class="phon">/daɪˈrekt/</span></div> <div class="phons_n_am" geo="n_am" htag="div" hclass="phons_n_am" wd="direct"><div class="sound audio_play_button pron-us icon-audio" data-src-mp3="https://www.oxfordlearnersdictionaries.com/media/english/us_pron/d/dir/direc/direct__us_1_rr.mp3" data-src-ogg="https://www.oxfordlearnersdictionaries.com/media/english/us_pron_ogg/d/dir/direc/direct__us_1_rr.ogg" title="direct pronunciation
                    American" style="cursor: pointer" valign="top">&nbsp;</div><span class="phon">/dəˈrekt/</span><span class="sep">,</span> <div class="sound audio_play_button pron-us icon-audio" data-src-mp3="https://www.oxfordlearnersdictionaries.com/media/english/us_pron/d/dir/direc/direct__us_2_rr.mp3" data-src-ogg="https://www.oxfordlearnersdictionaries.com/media/english/us_pron_ogg/d/dir/direc/direct__us_2_rr.ogg" title="direct pronunciation
                    American" style="cursor: pointer" valign="top">&nbsp;</div><span class="phon">/daɪˈrekt/</span></div></span></div>`

	levelExtractor := func(node *html.Node) (string, error) {
		for _, a := range node.Attr {
			if a.Key == "href" {
				ps := strings.Split(a.Val, "&")
				ls := ps[len(ps)-1]
				if strings.HasPrefix(ls, "level=") {
					return strings.Replace(ls, "level=", "", 1), nil
				}
			}
		}
		return "", nil
	}

	type Phonetic struct {
		Audio      []string `select:".sound" extract:"@data-src-mp3"`
		Transcript []string `select:".phon" extract:"text"`
	}

	type Phonetics struct {
		UK Phonetic `select:".phons_br"`
		US Phonetic `select:".phons_n_am"`
	}

	type Top struct {
		Term      string    `select:".headword" extract:"text"`
		Pos       string    `select:".pos" extract:"text"`
		Level     string    `select:".symbols a" extract:"*level"`
		Phonetics Phonetics `select:".phonetics"`
	}

	cfgs := []ScrapeCfg{
		{
			CaseName:   "scrape dictionary info",
			extractors: map[string]Extractor{"*level": levelExtractor},
			doc:        getDoc(htmldata),
			o:          &Top{},
			exp: &Top{
				Term:  "direct",
				Pos:   "adjective",
				Level: "a2",
				Phonetics: Phonetics{
					UK: Phonetic{
						Audio:      []string{"https://www.oxfordlearnersdictionaries.com/media/english/uk_pron/d/dir/direc/direct__gb_7.mp3", "https://www.oxfordlearnersdictionaries.com/media/english/uk_pron/d/dir/direc/direct__gb_8.mp3"},
						Transcript: []string{"/dəˈrekt/", "/daɪˈrekt/"},
					},
					US: Phonetic{
						Audio:      []string{"https://www.oxfordlearnersdictionaries.com/media/english/us_pron/d/dir/direc/direct__us_1_rr.mp3", "https://www.oxfordlearnersdictionaries.com/media/english/us_pron/d/dir/direc/direct__us_2_rr.mp3"},
						Transcript: []string{"/dəˈrekt/", "/daɪˈrekt/"},
					},
				},
			},
		},
	}
	tab.RunWithCfgs(t, cfgs, test)
}

func TestScraper_Scrape_ScrapeSlice(t *testing.T) {
	type S struct {
		Name string `extract:"@class"`
	}
	cfgs := []ScrapeCfg{
		{
			CaseName: "slice of strings",
			doc:      getDoc(`<div id="top"><div class="con">golang</div><div class="noc"></div></div>`),
			o:        &[]string{},
			selector: "#top > div",
			extract:  "@class",
			exp:      &[]string{"con", "noc"},
		},
		{
			CaseName: "slice of structs",
			doc:      getDoc(`<div id="top"><div class="con">golang</div><div class="noc"></div></div>`),
			o:        &[]S{},
			selector: "#top > div",
			exp:      &[]S{{"con"}, {"noc"}},
		},
	}
	tab.RunWithCfgs(t, cfgs, test)
}

func TestScraper_Scrape_ScrapeStruct(t *testing.T) {
	type Ex1 struct {
		Name  string `extract:"@class"`
		Value string `extract:"TEXT"`
	}
	type Ex2 struct {
		Name  string `extract:"@class"`
		Value string `extract:"text"`
	}

	cfgs := []ScrapeCfg{
		{
			CaseName: "not found",
			strict:   true,
			doc:      getDoc(`<div id="top"><div class="con">golang</div></div>`),
			o:        &Ex1{},
			selector: ".play",
			exp:      &Ex1{},
			eErr:     GetNotFoundErr(".play"),
		},
		{
			CaseName: "field extract tag err",
			strict:   true,
			doc:      getDoc(`<div id="top"><div class="con">golang</div></div>`),
			o:        &Ex1{},
			selector: "#top > .con",
			exp:      &Ex1{Name: "con"},
			eErr:     GetExtractErr("TEXT"),
		},
		{
			CaseName: "normal execution",
			doc:      getDoc(`<div id="top"><div class="con">golang</div></div>`),
			o:        &Ex2{},
			selector: "#top > .con",
			exp:      &Ex2{Name: "con", Value: "golang"},
		},
	}

	tab.RunWithCfgs(t, cfgs, test)
}

func TestScraper_Scrape_ScrapeString(t *testing.T) {
	s1, s2, s3 := "", "", "golang"
	cfgs := []ScrapeCfg{
		{
			CaseName: "not strict mod",
			doc:      getDoc(""),
			selector: "p",
			o:        &s1,
			exp:      &s2,
		},
		{
			CaseName: "strict mod",
			strict:   true,
			doc:      getDoc(`<div class="con">Text</div>`),
			o:        &s1,
			exp:      &s2,
			selector: ".cont",
			eErr:     GetNotFoundErr(".cont"),
		},
		{
			CaseName: "empty extract err",
			doc:      getDoc(`<div class="con">Text</div>`),
			o:        &s1,
			selector: ".con",
			exp:      &s2,
			eErr:     GetExtractErr(""),
		},
		{
			CaseName: "extractor err",
			doc:      getDoc(`<div class="con">Text</div>`),
			o:        &s1,
			exp:      &s2,
			selector: ".con",
			extract:  "@href",
			eErr:     WrapExtractErr(".con", GetAttributeNotFoundErr("href")),
		},
		{
			CaseName: "inner text",
			doc:      getDoc(`<div id="top"><div class="con">golang</div></div>`),
			o:        &s1,
			selector: "#top div",
			extract:  TextExtractTag,
			exp:      &s3,
		},
		{
			CaseName: "attribute",
			doc:      getDoc(`<div id="top"><div class="con" data="golang">text</div></div>`),
			o:        &s1,
			selector: "#top > div",
			extract:  "@data",
			exp:      &s3,
		},
	}
	tab.RunWithCfgs(t, cfgs, test)
}

func TestScraper_Scrape_CommonErrors(t *testing.T) {
	cfgs := []ScrapeCfg{
		{
			CaseName: "nil errors",
			eErr:     errors.Join(GetNilErr("doc"), GetNilErr("o")),
		},
		{
			CaseName: "not a pointer",
			doc:      getDoc(""),
			o:        "",
			exp:      "",
			eErr:     GetKindErr("string", "ptr", "string"),
		},
		{
			CaseName: "invalid kind",
			doc:      getDoc(""),
			o:        &map[int]int{},
			exp:      &map[int]int{},
			eErr:     GetMultiKindErr("map[int]int", []any{"string", "slice", "struct"}, "map"),
		},
	}
	tab.RunWithCfgs(t, cfgs, test)
}

func getDoc(data string) *goquery.Document {
	r := bytes.NewBufferString(data)
	doc, _ := goquery.NewDocumentFromReader(r)
	return doc
}
