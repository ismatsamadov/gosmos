package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

// localPath reports whether next is a path on this site, which a redirect may send a browser to.
func localPath(next string) bool {
	_, err := url.Parse(next)
	return err == nil && strings.HasPrefix(next, "/") && !strings.HasPrefix(next, "//") && !strings.Contains(next, `\`)
}

// goNext sends a browser that followed a link in one of Yusra's emails on to the page next names, and logs the click; a next that isn't a path on this site goes to the guide.
func goNext(l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next := r.URL.Query().Get("next")
		if !localPath(next) {
			next = "/guide/"
		}
		l.Printf("%s go %s", requestIDFrom(r.Context()), next)
		http.Redirect(w, r, next, http.StatusSeeOther)
	})
}
