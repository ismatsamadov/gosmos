package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

// articleSources are the articles, in the order the app lists them: where each file is, under
// $GOROOT or the golang/website checkout, and where go.dev publishes it, which relative links
// in it start from.
var articleSources = []struct {
	slug, file, page string
	website          bool
}{
	{"spec", "doc/go_spec.html", "/ref/spec", false},
	{"mem", "doc/go_mem.html", "/ref/mem", false},
	{"effective-go", "_content/doc/effective_go.html", "/doc/effective_go", true},
	{"faq", "_content/doc/faq.md", "/doc/faq", true},
	{"code", "_content/doc/code.html", "/doc/code", true},
}

// An article is a document read and split into its sections, ready to be written as Markdown.
type article struct {
	slug, title, subtitle string
	base                  *url.URL // the page on go.dev
	lead                  []*node  // what comes before the first heading
	sections              []*section
	ids                   map[string]string // each element id → the anchor of its section; "" for the lead
}

type section struct {
	anchor, title string
	level         int
	nodes         []*node
}

// loadArticles reads the five articles.
func loadArticles(goroot, website string) ([]*article, error) {
	var out []*article
	for _, s := range articleSources {
		root, file := goroot, filepath.FromSlash(s.file)
		if s.website {
			root = website
		}
		a, err := loadArticle(filepath.Join(root, file), filepath.Join(website, "_content"), s.page)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", s.slug, err)
		}
		a.slug = s.slug
		out = append(out, a)
	}
	return out, nil
}

// loadArticle reads a go.dev document as the site builds it: its metadata (a JSON comment at
// the top of HTML, YAML front matter in Markdown), then its template if it asks for one, then
// goldmark for Markdown. content is the website's _content folder, which {{code}} reads from.
func loadArticle(file, content, page string) (*article, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	src := string(data)
	var meta struct {
		Title, Subtitle string
		Template        bool
	}
	if strings.HasSuffix(file, ".md") {
		front, body, ok := strings.Cut(strings.TrimPrefix(src, "---\n"), "\n---\n")
		if !ok || !strings.HasPrefix(src, "---\n") {
			return nil, fmt.Errorf("%s has no front matter", file)
		}
		src = body
		for line := range strings.SplitSeq(front, "\n") {
			key, value, _ := strings.Cut(line, ":")
			value = strings.Trim(strings.TrimSpace(value), `"'`)
			switch strings.TrimSpace(key) {
			case "title":
				meta.Title = value
			case "subtitle":
				meta.Subtitle = value
			case "template":
				meta.Template = value == "true"
			}
		}
	} else if rest, ok := strings.CutPrefix(strings.TrimLeft(src, " \t\r\n"), "<!--{"); ok {
		js, body, ok := strings.Cut(rest, "}-->")
		if !ok {
			return nil, fmt.Errorf("%s: its metadata comment is not closed", file)
		}
		if err := json.Unmarshal([]byte("{"+js+"}"), &meta); err != nil {
			return nil, fmt.Errorf("%s: its metadata: %w", file, err)
		}
		src = body
	}
	if meta.Title == "" {
		return nil, fmt.Errorf("%s has no title", file)
	}
	if meta.Template {
		if src, err = runTemplate(src, content, path.Dir(page)); err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
	}
	if strings.HasSuffix(file, ".md") {
		var b bytes.Buffer
		md := goldmark.New(
			goldmark.WithExtensions(extension.Table, extension.DefinitionList),
			goldmark.WithParserOptions(parser.WithHeadingAttribute(), parser.WithAutoHeadingID()),
			goldmark.WithRendererOptions(gmhtml.WithUnsafe()),
		)
		if err := md.Convert([]byte(expandTabs(src)), &b); err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
		src = b.String()
	}
	a := &article{title: meta.Title, subtitle: meta.Subtitle, base: &url.URL{Scheme: "https", Host: "go.dev", Path: page}}
	a.split(parseHTML(src))
	return a, nil
}

// split divides a document into its lead and its sections, each heading starting one, and
// notes which section each element id is in.
func (a *article) split(doc *node) {
	a.ids = map[string]string{}
	used := map[string]bool{}
	var walk func(nodes []*node)
	walk = func(nodes []*node) {
		for _, n := range nodes {
			if level := headingLevel[n.tag]; level > 0 {
				title := strings.TrimSpace(reSpaces.ReplaceAllString(plain(n), " "))
				anchor := n.attr["id"]
				if anchor == "" || used[anchor] {
					anchor = uniqueAnchor(title, used)
				}
				used[anchor] = true
				a.ids[anchor] = anchor
				if id := n.attr["id"]; id != "" && a.ids[id] == "" {
					a.ids[id] = anchor
				}
				a.sections = append(a.sections, &section{anchor: anchor, title: title, level: level})
				continue
			}
			if containerTags[n.tag] && !n.hasClass("note") {
				walk(n.children)
				continue
			}
			current := ""
			if k := len(a.sections); k > 0 {
				a.sections[k-1].nodes = append(a.sections[k-1].nodes, n)
				current = a.sections[k-1].anchor
			} else {
				a.lead = append(a.lead, n)
			}
			eachID(n, func(id string) {
				if _, seen := a.ids[id]; !seen {
					a.ids[id] = current
				}
			})
		}
	}
	walk(doc.children)
}

// reProduction finds the names a grammar block defines, which go.dev makes anchors of.
var reProduction = regexp.MustCompile(`(?m)^\s*([A-Za-z_]\w*)\s*=`)

// eachID calls fn with the id of n and of every element inside it, the name of every <a name>,
// and every production a grammar block defines, as go.dev gives the spec's productions ids.
func eachID(n *node, fn func(string)) {
	if id := n.attr["id"]; id != "" {
		fn(id)
	}
	if n.tag == "a" && n.attr["name"] != "" {
		fn(n.attr["name"])
	}
	if n.tag == "pre" && n.hasClass("ebnf") {
		for _, m := range reProduction.FindAllStringSubmatch(plain(n), -1) {
			fn(m[1])
		}
	}
	for _, ch := range n.children {
		eachID(ch, fn)
	}
}

var reNotAnchor = regexp.MustCompile(`[^A-Za-z0-9]+`)

// uniqueAnchor makes an anchor from a heading's title for a heading without an id of its own.
func uniqueAnchor(title string, used map[string]bool) string {
	base := strings.Trim(reNotAnchor.ReplaceAllString(title, "_"), "_")
	if base == "" {
		base = "section"
	}
	anchor := base
	for i := 1; used[anchor]; i++ {
		anchor = fmt.Sprintf("%s_%d", base, i)
	}
	return anchor
}

// runTemplate runs a page's template as go.dev does, with the two functions the articles use:
// raw, which writes its argument as it is, and code, which shows lines of a file in content.
func runTemplate(src, content, dir string) (string, error) {
	t, err := template.New("page").Funcs(template.FuncMap{
		"raw": func(s any) string { return fmt.Sprint(s) },
		"code": func(file string, args ...any) (string, error) {
			if !path.IsAbs(file) {
				file = path.Join(dir, file)
			}
			data, err := os.ReadFile(filepath.Join(content, filepath.FromSlash(path.Clean(file))))
			if err != nil {
				return "", err
			}
			text, err := codeLines(string(data), args)
			if err != nil {
				return "", fmt.Errorf("code %s: %w", file, err)
			}
			return "<pre>" + html.EscapeString(text) + "</pre>\n", nil
		},
	}).Parse(src)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	if err := t.Execute(&b, nil); err != nil {
		return "", err
	}
	return b.String(), nil
}

// codeLines picks the lines of a file that go.dev's {{code}} shows: all of them, the line one
// pattern matches, or the lines from the one the first pattern matches to the one the second
// matches. A pattern is /regexp/, $ for the last line, or a line number. Lines ending in OMIT
// are left out, and tabs become four spaces.
func codeLines(text string, args []any) (string, error) {
	if n := len(args); n > 0 { // the lines to highlight, which a code block cannot show
		if s, ok := args[n-1].(string); ok && strings.HasPrefix(s, "HL") {
			args = args[:n-1]
		}
	}
	lines := strings.SplitAfter(text, "\n")
	find := func(start int, arg any) (int, error) {
		switch a := arg.(type) {
		case int:
			if a < 1 || a > len(lines) {
				return 0, fmt.Errorf("line %d is out of range", a)
			}
			return a, nil
		case string:
			if a == "$" {
				return len(lines), nil
			}
			if len(a) > 2 && a[0] == '/' && a[len(a)-1] == '/' {
				re, err := regexp.Compile("(?m)" + a[1:len(a)-1])
				if err != nil {
					return 0, err
				}
				for i := start; i < len(lines); i++ {
					if re.MatchString(lines[i]) {
						return i + 1, nil
					}
				}
				return 0, fmt.Errorf("no line matches %s", a)
			}
		}
		return 0, fmt.Errorf("unknown pattern %v", arg)
	}
	from, to := 1, len(lines)
	var err error
	switch len(args) {
	case 0:
	case 1:
		if from, err = find(0, args[0]); err == nil {
			to = from
		}
	case 2:
		if from, err = find(0, args[0]); err == nil {
			to, err = find(from, args[1])
		}
	default:
		err = fmt.Errorf("%d patterns, want at most 2", len(args))
	}
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for _, line := range lines[from-1 : to] {
		if !strings.HasSuffix(line, "OMIT\n") {
			b.WriteString(line)
		}
	}
	out := strings.Trim(b.String(), "\n")
	if out != "" {
		out += "\n"
	}
	return strings.ReplaceAll(out, "\t", "    "), nil
}

// expandTabs turns tabs into spaces up to four-column tab stops, as go.dev does before goldmark
// reads a page.
func expandTabs(s string) string {
	var b strings.Builder
	col := 0
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		s = s[size:]
		switch r {
		case '\n':
			b.WriteByte('\n')
			col = 0
		case '\t':
			b.WriteByte(' ')
			for col++; col%4 != 0; col++ {
				b.WriteByte(' ')
			}
		default:
			b.WriteRune(r)
			col++
		}
	}
	return b.String()
}
