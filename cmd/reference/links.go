package main

import (
	"net/url"
	"path"
	"strings"
)

// A linker knows what the reference holds, and turns the links of the sources into links to
// it: godoc: for packages and symbols, godoc-article: for articles and their sections. A link
// to anything else becomes its text, since nothing in the app may lead to the web.
type linker struct {
	symbols  map[string]map[string]string // import path → Name or Type.Name → kind
	articles map[string]*article          // by slug
}

// goDev is where relative links in doc comments lead.
var goDev = &url.URL{Scheme: "https", Host: "go.dev", Path: "/"}

// articlePaths are where the articles live on go.dev and golang.org, now and before.
var articlePaths = map[string]string{
	"/ref/spec":         "spec",
	"/doc/go_spec":      "spec",
	"/ref/mem":          "mem",
	"/doc/go_mem":       "mem",
	"/doc/effective_go": "effective-go",
	"/doc/faq":          "faq",
	"/doc/go_faq":       "faq",
	"/doc/code":         "code",
}

// target returns the reference's URL for href, a link on the go.dev page at base, or "" when it
// leads outside the reference.
func (l *linker) target(base *url.URL, href string) string {
	u, err := base.Parse(strings.TrimSpace(href))
	if err != nil || u.Scheme != "http" && u.Scheme != "https" {
		return ""
	}
	p := strings.TrimSuffix(strings.TrimSuffix(path.Clean("/"+u.Path), "/"), ".html")
	switch strings.TrimPrefix(u.Host, "www.") {
	case "go.dev", "golang.org", "tip.golang.org", "golang.google.cn":
		if slug, ok := articlePaths[p]; ok {
			return l.articleURL(slug, u.Fragment)
		}
		if pkg, ok := strings.CutPrefix(p, "/pkg/"); ok {
			return l.packageURL(pkg, u.Fragment)
		}
	case "pkg.go.dev":
		pkg, _, _ := strings.Cut(strings.TrimPrefix(p, "/"), "@") // pkg.go.dev/strings@go1.22
		return l.packageURL(pkg, u.Fragment)
	}
	return ""
}

// articleURL links to the article's section that holds the element with the id anchor, or to
// the article when there is no such element.
func (l *linker) articleURL(slug, anchor string) string {
	a := l.articles[slug]
	if a == nil {
		return ""
	}
	if section := a.ids[anchor]; anchor != "" && section != "" {
		return "godoc-article:" + slug + "#" + section
	}
	return "godoc-article:" + slug
}

// packageURL links to the package at importPath, or to its symbol when anchor, pkg.go.dev's
// fragment, names one: Name, or Type.Name for a method, field or interface method.
func (l *linker) packageURL(importPath, anchor string) string {
	recv, name, ok := strings.Cut(anchor, ".")
	if !ok {
		recv, name = "", anchor
	}
	return l.symbolURL(importPath, recv, name)
}

// symbolURL links to a symbol of the package at importPath: recv.name, or name when recv is
// "". A symbol the package does not have links to the package, and a package the reference
// does not have links nowhere.
func (l *linker) symbolURL(importPath, recv, name string) string {
	syms, ok := l.symbols[importPath]
	switch {
	case !ok:
		return ""
	case recv != "" && (syms[recv+"."+name] == "method" || syms[recv] == "type"):
		return "godoc:" + importPath + "." + recv + "." + name
	case recv == "" && name != "" && syms[name] != "":
		return "godoc:" + importPath + "." + name
	}
	return "godoc:" + importPath
}
