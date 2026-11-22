package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

// fetch sends a GET request to url and returns the answer's status and body.
func fetch(url string) (status, body string, err error) {
	_ = http.Get
	_ = io.ReadAll
	return "", "", nil
}

func main() {
	srv := httptest.NewServer(standIn())
	defer srv.Close()
	for _, code := range []string{"PEL", "ZZZ"} {
		status, body, err := fetch(srv.URL + "/stations/" + code)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Println(status, strings.TrimSpace(body))
	}
}
