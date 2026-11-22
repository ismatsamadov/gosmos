---
slug: tell-them-when
title: Tell Them When
type: code_tests
position: 3
difficulty: 8
topics: http, time
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Yusra's TL;DR from Brinewatch: "We'd gladly wait, if your service told us how long." Today a refused request says nothing, so their client tries again at once.

`limit` runs after `authenticate`, which leaves the customer in the request's context. A request over the limit gets `429 Too Many Requests` and `Retry-After` in whole seconds, and every answer carries `X-RateLimit-Limit` and `X-RateLimit-Remaining`, so a program can slow down before it is refused:

```go
wait := 2 * time.Second
fmt.Println(wait)                    // 2s
fmt.Println(int(wait / time.Second)) // 2
```

## Your task

`limit` lets everything through. Take the customer with `customerFrom`; without one, answer through `writeError` with an error of your own, which makes a 500. Call `l.Allow` with its ID, set `X-RateLimit-Limit` to `10` and `X-RateLimit-Remaining` to what is left, then refuse with `Retry-After` and `errRateLimited`, or call `next`.

> [!NOTE]
> A request without a customer means the routes are wired wrong: a 500 says so, where letting it through would leave it unlimited.
