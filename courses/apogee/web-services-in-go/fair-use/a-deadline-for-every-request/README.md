---
slug: a-deadline-for-every-request
title: A Deadline for Every Request
type: code_tests
position: 6
difficulty: 9
topics: context, errors, http
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
During the flood, Cirrine's bookings waited half a minute for the store, and nobody answered them until it was over. A request that can't be answered in time should hear so, and soon.

`context.WithTimeout` makes a context that ends after a while, and whatever watches it, the store included, gives up then with `context.DeadlineExceeded`. Middleware can give every request such a context through `r.WithContext`, and must cancel it when the handler returns, so its timer stops:

```go
ctx, cancel := context.WithTimeout(
	context.Background(), 50*time.Millisecond)
defer cancel()
<-ctx.Done()
fmt.Println(ctx.Err())
// context deadline exceeded
```

## Your task

1. `withDeadline` returns `next` unchanged. When `d` is above zero, call `next` with a request whose context ends after `d`, and cancel that context when `next` returns.
2. In `writeError`, answer `context.DeadlineExceeded` with 503 `unavailable` `request took too long`, and `context.Canceled` with 503 `unavailable` `request canceled`.
