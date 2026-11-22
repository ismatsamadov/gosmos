package main

import "strings"

// greeting is the first line the watch port sends a watcher: the dishes dishd knows.
func greeting(names []string) string {
	return "dishd watching " + strings.Join(names, " ") + "\n"
}
