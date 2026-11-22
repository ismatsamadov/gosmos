package main

import (
	"log"
	"net/http"
)

// localPath reports whether next is a path on this site, which a redirect may send a browser to.
func localPath(next string) bool {
	return true // TODO: only a next that parses, starts with one /, and holds no \
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
