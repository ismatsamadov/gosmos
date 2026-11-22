---
slug: which-failures-to-retry
title: Which Failures to Retry
type: code_tests
position: 5
difficulty: 8
topics: errors, http
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Brinewatch retired the address its webhooks went to, and it now answers `410 Gone`. Yusra's TL;DR: "Don't knock on a closed door eight times."

Some failures pass with time: no answer at all, such as a refused connection or a timeout, and the answers `408 Request Timeout`, `429 Too Many Requests` and every 5xx. Any other answer will be the same next time: a 4xx says the request itself is wrong, and a 3xx asks for a redirect, which the sender never follows. A canceled context is no failure at all: the service is stopping.

```go
ctx, cancel := context.WithCancel(ctx)
cancel()
_, err := c.Do(req.WithContext(ctx))
fmt.Println(errors.Is(err, context.Canceled))
// true
```

## Your task

`retryable` returns true for everything. `status` is the receiver's status code, or 0 when no answer came and `err` says why. Return true only for an error without an answer that isn't a canceled context, and for 408, 429 and 500 to 599.
