---
slug: every-request-gets-an-id
title: Every Request Gets an ID
type: code_tests
position: 1
difficulty: 7
topics: context, http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Cirrine writes that "one of this morning's requests" failed, and Yusra can't tell which one. Nkem's plan: every request gets an ID, in its answer's `X-Request-Id` header and in every error body, for customers to quote.

A request carries a context. `context.WithValue` returns a new context holding one more value under a key, and `r.WithContext(ctx)` a copy of the request carrying it. A key of your own unexported type can't clash with anyone else's:

```go
type stationKey struct{}

ctx := context.WithValue(context.Background(),
	stationKey{}, "SQT")
code, ok := ctx.Value(stationKey{}).(string)
fmt.Println(code, ok) // SQT true
```

## Your task

1. Declare an unexported key type. In `requestID`, put the ID in the context under it, and call `next` with `r.WithContext`.
2. Make `requestIDFrom` return that ID, or `""` when there is none. Nkem has made `writeError` copy it into `request_id`.

> [!NOTE]
> A value stored under the plain string `"request_id"` isn't yours: `requestIDFrom` must not find it.
