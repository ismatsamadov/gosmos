---
slug: a-panic-gets-an-answer
title: A Panic Gets an Answer
type: code_tests
position: 5
difficulty: 8
topics: defer, http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Monday's bug: Cirrine's program asked for a pass by its pass ID with a digit missing, `cirrine-1/YST/1136`, and `getPass` panicked. Cirrine got a dropped connection: no answer, nothing to quote. Bugs will slip through again; Nkem wants the next one to end in an answer.

A deferred function runs even while a panic unwinds, and inside it `recover()` stops the panic and returns its value. One value must go on: `http.ErrAbortHandler`, which a handler panics with on purpose to cut the connection.

```go
func safely(f func()) {
	defer func() {
		if v := recover(); v != nil {
			fmt.Println("recovered:", v)
		}
	}()
	f()
}
```

## Your task

`recoverPanics` only calls `next`. Defer a function that recovers a panic, then:

1. panics again with the value, untouched, if it is `http.ErrAbortHandler`;
2. otherwise logs this line and answers through `writeError` with an error of your own, which becomes a 500:

```text
req-000003 panic: <the panic's value>
```
