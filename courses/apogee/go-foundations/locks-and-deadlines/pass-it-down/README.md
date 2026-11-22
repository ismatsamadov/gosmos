---
slug: pass-it-down
title: Pass It Down
type: code
position: 5
difficulty: 7
topics: context, functions
created: 2026-11-01T07:36:47Z
modified: 2026-11-06T16:51:27Z
---
On Tuesday a stuck uplink held La Grande at Salar Quieto for an hour. The pass had a deadline, but one function in the chain started a fresh context, so `transmit` never heard about it.

A context only works if it reaches every call that waits. By convention it is the first parameter, named `ctx`, and each function passes its own `ctx` on instead of starting over with `context.Background()`.

```go
// fetchElements downloads fresh orbit elements,
// giving up when ctx ends.
func fetchElements(
	ctx context.Context, satID string,
) ([]byte, error) {
	return download(ctx, "elements/"+satID)
}
```

## Your task

Give `sendCommand` a `ctx context.Context` first parameter and pass it to `transmit`. Then pass `runPass`'s `ctx` to `sendCommand`.

> [!TIP]
> Pass `ctx` to each call; don't store it in a struct field, where it would outlive that call.
