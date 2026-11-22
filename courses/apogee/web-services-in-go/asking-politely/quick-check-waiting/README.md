---
slug: quick-check-waiting
title: Quick Check: Waiting
type: quiz
position: 2
difficulty: 5
topics: context, http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Tomás wrote the weather service in 2024, and he won't hear it called slow: "It was never slow. It was patient!" The poll's old client was patient too. Before Ingrid trusts the new one, three questions about this code, which calls a server that has stalled:

```go
c := &http.Client{Timeout: 2 * time.Second}
ctx, cancel := context.WithTimeout(
	context.Background(), 500*time.Millisecond)
defer cancel()
req, _ := http.NewRequestWithContext(ctx,
	http.MethodGet, url, nil)
resp, err := c.Do(req)
```

## Your task

Answer three questions: how long `http.Get` waits, what a client's `Timeout` covers, and when `c.Do` returns in the code above.
