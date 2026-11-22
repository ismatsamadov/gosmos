---
slug: a-handler-around-a-handler
title: A Handler Around a Handler
type: code_tests
position: 0
difficulty: 7
topics: functions, http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Cirrine canceled booking `bk-4511`, but their office proxy, a machine that relays requests and keeps copies of answers, kept the old answer to `GET /bookings/bk-4511`: for a day, their program showed a canceled booking.

Nkem wants the fix on every answer, so it goes in middleware: a function that takes a handler and returns a new one, which does its part, then calls the handler it wraps. `Cache-Control: no-store` asks every cache on the way, such as that proxy, to keep no copy.

```go
func english(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter,
		r *http.Request) {
		w.Header().Set("Content-Language", "en")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}
```

## Your task

`noStore` returns `next` unchanged. Make it return a handler that sets `Cache-Control: no-store` on `w`, then calls `next`, so every answer carries the header, the mux's own 404s and 405s too.

> [!WARNING]
> Set the header before calling `next`: once `next` has written the answer, a header set afterward changes nothing.
