---
slug: what-an-answer-says
title: What an Answer Says
type: code
position: 4
difficulty: 4
topics: http
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Time to hear an answer. Nkem's stand-in, a few lines of Go in `standin.go`, answers like the booking service will, about stations only, and you can ask it questions with no network at all.

The stand-in is a handler: a value with a `ServeHTTP(w, r)` method, which reads the request `r` and writes its answer to `w`. `httptest.NewRequest` makes an `r`, and `httptest.NewRecorder` makes a `w` that keeps everything written to it:

```go
rec := httptest.NewRecorder()
req := httptest.NewRequest(http.MethodGet,
	"/stars", nil)
http.NotFoundHandler().ServeHTTP(rec, req)
fmt.Println(rec.Code) // 404
```

## Your task

`main` makes a request for `/stations/PEL` and a recorder, but never hands them to the stand-in. Replace the `TODO` line with a call to the `ServeHTTP` method of `standIn()`, so `printResponse` prints:

```text
HTTP/1.1 200 OK
Content-Type: application/json

{"code":"PEL","name":"Port Ellery"}
```

> [!NOTE]
> Run the starter first: a recorder nobody wrote to still says `200 OK`, the status every answer has until a handler sets another.
