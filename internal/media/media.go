// Package media finds what the app never shows in a lesson's or a story's Markdown: pictures
// and HTML. A library holds none of them, since the app shows lessons as text, offline, so the
// authoring tools refuse a text that has any.
//
// Code is never looked into: fenced and indented code blocks and code spans are text, because
// a lesson may teach HTML or Markdown, and its examples (an <img> in a test page, an <iframe>
// in a security example, the image syntax for a README badge) are code.
package media

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Found counts what Find found.
type Found struct {
	Images int // pictures: Markdown's image syntax, a video file put in it too
	HTML   int // HTML blocks and inline HTML: <img>, <video>, <iframe>, comments and any other tag
}

// Total is all Find found.
func (f Found) Total() int { return f.Images + f.HTML }

// String says what was found, such as "1 picture and 2 pieces of HTML".
func (f Found) String() string {
	var parts []string
	if f.Images > 0 {
		parts = append(parts, count(f.Images, "picture", "pictures"))
	}
	if f.HTML > 0 {
		parts = append(parts, count(f.HTML, "piece of HTML", "pieces of HTML"))
	}
	if len(parts) == 0 {
		return "nothing"
	}
	return strings.Join(parts, " and ")
}

func count(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}
	return fmt.Sprintf("%d %s", n, many)
}

var mdParser = goldmark.New().Parser()

// Find counts the pictures and the HTML in md, outside code, as CommonMark reads it.
func Find(md string) Found {
	var f Found
	doc := mdParser.Parse(text.NewReader([]byte(md)))
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch n.(type) {
			case *ast.Image:
				f.Images++
			case *ast.HTMLBlock, *ast.RawHTML:
				f.HTML++
			}
		}
		return ast.WalkContinue, nil
	})
	return f
}

// OneParagraph reports whether md is a single Markdown paragraph: text that may use inline
// Markdown, such as code spans, emphasis and links, but no headings, lists, quotes, code
// blocks, HTML blocks or thematic breaks.
func OneParagraph(md string) bool {
	doc := mdParser.Parse(text.NewReader([]byte(md)))
	first := doc.FirstChild()
	return first != nil && first.NextSibling() == nil && first.Kind() == ast.KindParagraph
}
