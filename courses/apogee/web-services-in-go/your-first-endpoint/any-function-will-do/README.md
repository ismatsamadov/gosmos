---
slug: any-function-will-do
title: Any Function Will Do
type: code_tests
position: 1
difficulty: 5
topics: functions, http, methods
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Nkem reads `stationList` and asks one question: "A new type for every answer?" She would rather keep each handler as a method of one `Service`, which holds the stations now and everything else later.

`http.HandlerFunc` is a function type whose `ServeHTTP` method calls the function itself. Converting any function with a handler's two parameters to it makes a handler, and a method value such as `s.listStations` is such a function, one that remembers its `s`:

```go
hello := func(w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprintln(w, "Port Ellery standing by")
}
var h http.Handler = http.HandlerFunc(hello)
```

## Your task

`stationsHandler` returns `http.NotFoundHandler()`, so every request gets a 404. Return `s.listStations` as a handler instead, with `http.HandlerFunc`.

> [!NOTE]
> `http.HandlerFunc(hello)` is a conversion, not a call: `hello` runs later, once for every request.
