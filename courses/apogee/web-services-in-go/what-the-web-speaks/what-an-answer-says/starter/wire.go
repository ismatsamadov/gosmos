package main

import (
	"fmt"
	"maps"
	"net/http/httptest"
	"slices"
	"strings"
)

// printResponse prints the answer rec caught: the status line, the headers, a blank line, then the body.
func printResponse(rec *httptest.ResponseRecorder) {
	res := rec.Result()
	fmt.Println(res.Proto, res.Status)
	for _, name := range slices.Sorted(maps.Keys(res.Header)) {
		fmt.Printf("%s: %s\n", name, strings.Join(res.Header[name], ", "))
	}
	fmt.Println()
	fmt.Print(rec.Body.String())
}
