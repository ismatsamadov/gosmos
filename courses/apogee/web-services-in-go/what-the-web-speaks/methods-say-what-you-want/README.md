---
slug: methods-say-what-you-want
title: Methods Say What You Want
type: code
position: 1
difficulty: 4
topics: http
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Yusra forwards three of this morning's emails with a one-line summary on top: "TL;DR: one asks which passes are free, one books a pass, one cancels booking `bk-4494`." Programs will send the same three wishes, and the first word of each request says which it is.

That word is the method. `GET` reads and changes nothing, `POST` creates something new, such as a booking, and `DELETE` removes something. One path can mean two things: `GET /bookings/bk-4494` reads that booking, and `DELETE` on the same path cancels it.

```go
req, err := http.NewRequest(http.MethodGet,
	"http://127.0.0.1/stations/CWR", nil)
if err != nil {
	return err
}
fmt.Println(requestLine(req))
// GET /stations/CWR HTTP/1.1
```

## Your task

In `main`, all three requests use `http.MethodGet`. Give the booking `http.MethodPost` and the cancellation `http.MethodDelete`, so the program prints:

```text
GET /passes HTTP/1.1
POST /bookings HTTP/1.1
DELETE /bookings/bk-4494 HTTP/1.1
```
