package main

import (
	"fmt"
	"net/http"
)

func main() {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1/stations", nil)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	req.Header.Set("Accept", "application/json")
	printRequest(req)
}
