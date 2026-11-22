package main

import (
	"fmt"
	"net/http"
)

func main() {
	emails := []struct{ method, url string }{
		{http.MethodGet, "http://127.0.0.1/passes"},           // which passes are free?
		{http.MethodGet, "http://127.0.0.1/bookings"},         // book a pass
		{http.MethodGet, "http://127.0.0.1/bookings/bk-4494"}, // cancel booking bk-4494
	}
	for _, e := range emails {
		req, err := http.NewRequest(e.method, e.url, nil)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Println(requestLine(req))
	}
}
