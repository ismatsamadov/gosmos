---
slug: one-record-per-panic
title: One Record per Panic
type: code_tests
position: 8
difficulty: 7
topics: defer, logging
created: 2026-11-13T11:32:35Z
modified: 2026-11-13T11:32:35Z
---
Friday, 09:30, on the lab bench. Nkem gives uplinkd a pass ID that has lost its orbit number, `cirrine-1/LAN`, and preparing the pass panics. The guard around each pass keeps the station flying and writes the pass's failure line, but `debug.PrintStack` sends where it broke to standard error, as plain lines no search reads.

`debug.Stack`, from `runtime/debug`, returns the calling goroutine's stack as text in a `[]byte`. Called in the deferred function that recovered, it still shows where the panic happened. As a string in one attribute, the whole stack stays in one record, however many lines it has: the JSON handler writes each line break as `\n`.

```go
func where() string {
	return string(debug.Stack())
}
// goroutine 1 [running]:
// runtime/debug.Stack()
// ...
// main.where()
```

## Your task

In `guardPass`, give the failure line the stack, under the key `stack`, and stop printing it.

> [!WARNING]
> Given the `[]byte` itself, the JSON handler writes the stack in base64.
