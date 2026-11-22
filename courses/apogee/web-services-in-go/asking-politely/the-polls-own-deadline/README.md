---
slug: the-polls-own-deadline
title: The Poll's Own Deadline
type: code_tests
position: 1
difficulty: 5
topics: context, http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Ingrid's checklist has to be done ten minutes before AOS: "Checklist closes at AOS minus ten. Every call in it too." The poll hands `Reading` a context that ends at that moment, and the client's own two seconds may be longer than what is left.

A context carries a deadline, or a cancellation, from the caller to the work it starts. `http.NewRequestWithContext` ties a request to one: when the context ends, `Do` stops and returns its error, whatever the client's `Timeout` says.

```go
ctx, cancel := context.WithTimeout(
	context.Background(), 300*time.Millisecond)
defer cancel()
req, err := http.NewRequestWithContext(ctx,
	http.MethodGet, srv.URL+"/stations", nil)
```

## Your task

`Reading` builds its request with `http.NewRequest`, so the poll's deadline never reaches it. Build it with `http.NewRequestWithContext` and `ctx` instead.

> [!NOTE]
> Whichever ends first, the context or the client's `Timeout`, ends the call. An already canceled context ends it before it starts.
