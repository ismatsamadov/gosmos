---
slug: mirelas-review
title: Mirela's Review
type: code_tests
position: 9
difficulty: 9
topics: http, security, testing
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Friday 7 January. Mirela asks you to write her review as code, and has planted four broken copies of the service for it to catch: "Your review passes the service. Now show me the four it refuses."

A check earns its place when a broken service fails it. `review` gets the address of a running service and `advance`, which moves the service's clock forward. `calls.go` holds the beta customers' API keys, the private routes, and helpers that send a request and read the answer:

```go
status, msg := call(t, base, "GET /health", "")
// 200, and "" since it isn't an error
tok := token(t, base, cirrineKey)
list := bookingsOf(t, base, tok)
```

## Your task

Write `review` in `review.go`, with `t.Errorf` for every wrong answer:

1. Each of `privateRoutes`, sent with no credentials, gets 401.
2. Furrowlight's token lists only Furrowlight's bookings, and its API key sent as a `Bearer` token gets 401.
3. Brinewatch's token gets 404 for Cirrine's `bk-4510`, and after `advance(time.Hour)` Furrowlight's token gets `token expired`.
