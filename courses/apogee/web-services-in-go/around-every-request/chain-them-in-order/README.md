---
slug: chain-them-in-order
title: Chain Them in Order
type: code_tests
position: 6
difficulty: 7
topics: functions, http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Four middlewares wrapped by hand read inside out, with the mux buried in the middle. Nkem wants the service's wrapping written as one list that reads top to bottom.

A variadic parameter, `mws ...T`, takes any number of arguments and holds them as a slice. `chain(h, a, b, c)` must return `a(b(c(h)))`: the first one listed is outermost, so it sees a request first and the answer last.

```go
func total(base int, extra ...int) int {
	for _, e := range extra {
		base += e
	}
	return base
}
fmt.Println(total(10, 2, 3)) // 15
```

## Your task

`chain` returns `h` unwrapped. Make it wrap `h` in every middleware of `mws`, starting with the last one, and return the result, so that a request passes through them in the order they are listed.

> [!WARNING]
> Wrapping in list order makes the last one outermost: three tracing middlewares then record `c b a handler a b c`.
