---
slug: quick-check-middleware
title: Quick Check: Middleware
type: quiz
position: 2
difficulty: 7
topics: context, http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Before the service gets any more wrapping, Nkem has three questions, one at a time. The first is about this middleware, which the service calls once, when it starts, to wrap its mux:

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

Answer the three questions: when this code sets `Content-Language`, what `r.WithContext` does to `r`, and why the request ID's key has a type of its own.
