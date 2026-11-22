---
slug: a-type-that-answers
title: A Type That Answers
type: code_tests
position: 0
difficulty: 5
topics: http, interfaces, json
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Yusra's TL;DR for the week: "Cirrine's engineers start their client, and its first question is which stations exist." It is the first thing a customer will ever ask the booking service itself.

On the server's side, an answer comes from a handler: any value with a `ServeHTTP(w, r)` method, all that the `http.Handler` interface asks for. A type of yours satisfies it once it has that method, and what the method writes to `w` is the answer.

```go
type motto string

func (m motto) ServeHTTP(
	w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, string(m))
}

var h http.Handler = motto("Every pass counts")
```

## Your task

`stationList` has a `ServeHTTP` method that writes nothing. Make it encode the list to `w` with `json.NewEncoder`, so a request gets the stations as JSON in the list's order, and an empty list answers `[]`.

> [!WARNING]
> `fmt.Fprint(w, s)` writes Go's own printing of the slice, `[{YST Ystmark ...`, which no JSON client can read.
