---
slug: what-a-request-says
title: What a Request Says
type: code
position: 0
difficulty: 4
topics: http
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
From February, customers' programs will book passes from Apogee's own program, the booking service, instead of by email. Nkem sketches the first thing Furrowlight's program will ever send it: a request for the list of stations.

A request is lines of text. The first holds the method, the path and the version; each line after it is a header, a name and a value. `http.NewRequest` builds a request, and `req.Header.Set` adds a header:

```go
req, err := http.NewRequest(http.MethodGet,
	"http://127.0.0.1/health", nil)
if err != nil {
	return err
}
req.Header.Set("Accept-Language", "en")
```

## Your task

In `main`, make the request ask for `/stations`, and set its `Accept` header to `application/json`, which asks for the answer as JSON. The program should print:

```text
GET /stations HTTP/1.1
Host: 127.0.0.1
User-Agent: Go-http-client/1.1
Accept: application/json
```

> [!NOTE]
> Go writes `Host` and `User-Agent` itself. On the wire every line ends in two characters, `\r\n`, which `printRequest` prints as one line break.
