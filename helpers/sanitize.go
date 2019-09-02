package helpers

import (
	"bytes"
	"io"
	"reflect"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

func HTMLEscapeAll(this interface{}) {
	s := reflect.ValueOf(this).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.IsValid() && f.CanSet() {
			if f.Kind() == reflect.String {
				st := StripHtmlTags(strings.NewReader(f.String())).String()
				s.Field(i).SetString(st)
			}
		}
	}
}

// Parses the supplied HTML found in the io.reader
// It iterates over the the tags using the IsAllowedTag()
// method in order to find tags that are allowed.
func ParseHtml(r io.Reader) *bytes.Buffer {
	ti := html.NewTokenizer(r)
	b := bytes.NewBuffer(make([]byte, 0))
	for {
		tokenType := ti.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := ti.Token()
		// for now, strip all HTML attributes
		token.Attr = nil
		if IsAllowedTag(token) || (tokenType == html.TextToken) {
			b.WriteString(token.String())
		}
	}
	return b
}

// Parses the supplied HTML found in the io.reader
// It iterates over the the tags using the IsAllowedTag()
// method in order to find tags that are allowed.
func ParseCommentHtml(r io.Reader) *bytes.Buffer {
	ti := html.NewTokenizer(r)
	b := bytes.NewBuffer(make([]byte, 0))
	var breakTotal int
	for {
		tokenType := ti.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := ti.Token()
		// for now, strip all HTML attributes
		token.Attr = nil
		allow, thisTagIsBreak := IsAllowedCommentTag(token)

		if allow || (tokenType == html.TextToken) {
			if thisTagIsBreak {
				breakTotal = breakTotal + 1
			} else {
				breakTotal = 0
			}

			if breakTotal < 3 {
				b.WriteString(token.String())
			}

		}
	}
	return b
}

func StripHtmlTags(r io.Reader) *bytes.Buffer {
	ti := html.NewTokenizer(r)
	b := bytes.NewBuffer(make([]byte, 0))

	for {
		tokenType := ti.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := ti.Token()
		// for now, strip all HTML attributes
		token.Attr = nil
		if tokenType == html.TextToken {
			b.WriteString(token.String())
		}
	}
	return b
}

func ParseBreaksHtml(r io.Reader) *bytes.Buffer {
	ti := html.NewTokenizer(r)
	b := bytes.NewBuffer(make([]byte, 0))
	var breakTotal int
	for {
		tokenType := ti.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := ti.Token()
		thisTagIsBreak := IsBreakTag(token)
		if thisTagIsBreak {
			breakTotal = breakTotal + 1
		} else {
			breakTotal = 0
		}

		if breakTotal < 3 {
			b.WriteString(token.String())
		}

	}
	return b
}

func IsBreakTag(t html.Token) bool {
	if t.Data == "br" {
		return true
	}

	return false
}

func IsAllowedCommentTag(t html.Token) (bool, bool) {
	isBreak := false
	tags := map[string]bool{
		"strong": true,
		"em":     true,
		"i":      true,
		"ul":     true,
		"li":     true,
		"br":     true,
	}

	if t.Data == "br" {
		isBreak = true
	}

	if tags[t.Data] {
		return true, isBreak
	}

	return false, isBreak
}

func IsAllowedTag(t html.Token) bool {
	tags := map[string]bool{
		"strong": true,
		"p":      true,
		"em":     true,
		"i":      true,
		"h1":     true,
		"h2":     true,
		"h3":     true,
		"h4":     true,
		"ul":     true,
		"li":     true,
		"br":     true,
	}
	if tags[t.Data] {
		return true
	}
	return false
}

// allow a tiny subset of html elements with only the href attr on <a>
func Safe(s string) string {
	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()
	p.AllowAttrs("href").OnElements("a")
	p.AllowElements("strong")
	p.AllowElements("ul")
	p.AllowElements("li")
	p.AllowElements("p")
	p.AllowElements("em")
	p.AllowElements("br")
	p.AllowElements("del")
	p.AllowElements("i")
	s = html.UnescapeString(s)
	return p.Sanitize(s)
}

func SafeNoHref(s string) string {
	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()
	p.AllowElements("strong")
	p.AllowElements("ul")
	p.AllowElements("li")
	p.AllowElements("p")
	p.AllowElements("em")
	p.AllowElements("br")
	p.AllowElements("del")
	p.AllowElements("i")
	s = html.UnescapeString(s)
	return p.Sanitize(s)
}

// allow a tiny subset of html elements with only the href attr on <a>
func SafeTsAndCs(s string) string {
	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()
	p.AllowAttrs("href").OnElements("a")
	p.AllowElements("strong")
	p.AllowElements("ul")
	p.AllowElements("li")
	p.AllowElements("p")
	p.AllowElements("em")
	p.AllowElements("br")
	p.AllowElements("del")
	p.AllowElements("i")
	p.AllowElements("h1")
	p.AllowElements("h2")
	p.AllowElements("h3")
	s = html.UnescapeString(s)
	return p.Sanitize(s)
}

// take the bleve html results and make them safe for consumption
func BleveSafe(s string) string {
	p := bluemonday.NewPolicy()
	p.AllowElements("mark")
	s = html.UnescapeString(s)
	return p.Sanitize(s)
}

// strip it all, no html
func Strip(s string) string {
	p := bluemonday.StrictPolicy()
	return p.Sanitize(s)
}
