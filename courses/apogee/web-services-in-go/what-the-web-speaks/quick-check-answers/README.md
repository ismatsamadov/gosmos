---
slug: quick-check-answers
title: Quick Check: Answers
type: quiz
position: 7
difficulty: 5
topics: http
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Tomorrow the stand-in answers over a real connection for the first time. Before it does, Nkem has three questions about this code, which asks the stand-in about Salar Quieto:

```go
rec := httptest.NewRecorder()
req := httptest.NewRequest(http.MethodGet,
	"/stations/SQT", nil)
standIn().ServeHTTP(rec, req)
ct := rec.Header().Get("content-type")
```

## Your task

Answer the three questions: what a `201` says that a `200` doesn't, what `ct` holds once the code has run, and how the request reaches the stand-in.
