---
slug: the-flood-test
title: The Flood Test
type: code_tests
position: 7
difficulty: 9
topics: context, http, testing
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Friday 14 January: Nkem replays Brinewatch's flood against the service while Cirrine books its week. "What does Cirrine see while Brinewatch floods?"

The order of middleware decides what each one knows. `limit` reads the customer that `authenticate` puts in the context, so it must come after it, and then a request that `authenticate` refuses never spends anyone's tokens. `withDeadline` comes last of all, right around the mux, so it bounds the handlers' work and nothing else:

```go
h := chain(mux, first, second)
// h is first(second(mux)): a request
// meets first, then second, then mux
```

## Your task

1. In `auth.go`, make `private` list `authenticate` before `limit`.
2. In `middleware.go`, make `handler` list `withDeadline(s.Timeouts.Handler)` last, after `noStore`.

> [!NOTE]
> The test sends Brinewatch's 1,000 requests through `ServeHTTP` on a clock that stands still, so the whole flood takes a few milliseconds.
