package main

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
	"time"
)

// pages holds the Board's templates, built into the program.
//
//go:embed board.html script.html
var pages embed.FS

// pageData is what board.html is executed with. View is embedded, so the template
// says {{range .Passes}}.
type pageData struct {
	View
	Now time.Time // when the page is made, for its heading
}

// parsePage parses the Board's templates, once, for NewBoard.
func parsePage() (*template.Template, error) {
	t := template.New("board.html").Funcs(template.FuncMap{"asHTML": asHTML})
	return t.ParseFS(pages, "board.html", "script.html")
}

// asHTML marks s as trusted HTML, which a template prints as it is.
func asHTML(s string) template.HTML {
	return template.HTML(s)
}

// servePage writes the Board as an HTML page: the current view at the current time.
func (b *Board) servePage(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	if err := b.page.Execute(&buf, pageData{View: b.current(), Now: b.now()}); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}
