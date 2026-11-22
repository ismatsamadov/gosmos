---
slug: one-test-two-stores
title: One Test, Two Stores
type: code_tests
position: 3
difficulty: 8
topics: interfaces, testing
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
A store on disk is next, and Nkem wants it held to every promise the memory store keeps, by one test for any store. She has planted three copies of the store with one bug each, mutants for your test to catch.

A contract test is a function that takes a way to make what it checks, `open func() Store`. Each subtest opens a fresh store, so no check inherits leftovers:

```go
func checkCounter(t *testing.T,
	open func() Counter) {
	t.Run("starts at zero", func(t *testing.T) {
		if n := open().Count(); n != 0 {
			t.Errorf("Count() = %d, want 0", n)
		}
	})
}
```

## Your task

Write `checkStore` in `contract.go`, one subtest per check, with `t.Errorf` for each wrong answer:

1. The first two bookings get `bk-4501` and `bk-4502`, and `Get` returns what `Create` kept.
2. A second booking of the same pass returns `ErrConflict`; after `Delete`, `Get` returns `ErrNotFound`.
3. `List` orders by AOS: book the later pass first.
