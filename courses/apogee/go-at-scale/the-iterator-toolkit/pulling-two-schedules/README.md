---
slug: pulling-two-schedules
title: Pulling Two Schedules
type: code_tests
position: 5
difficulty: 9
topics: generics, iterators
created: 2026-11-11T22:52:57Z
modified: 2026-11-12T08:43:35Z
---
Wednesday 4 April. Since Monday, passplan 2 has planned beside the old planner and flown nothing: the shadow run. Ingrid wants the test Board to show it as two columns, the old planner's contacts beside the new one's, and a range loop walks one sequence at a time.

`iter.Pull` turns a sequence you range over into one you pull from. It returns two functions: `next` gives the next value and true, or the zero value and false once the sequence has ended, and `stop` ends it early. Until one of those happens, the sequence waits inside its loop, so call `stop` on every way out, with `defer`:

```go
next, stop := iter.Pull(slices.Values(codes))
defer stop()
if code, ok := next(); ok {
	fmt.Println("first:", code)
}
```

`Zip` yields its pairs as an `iter.Seq2`, the way `At` yielded each dish with its contact.

## Your task

`Zip` yields nothing. Make it range over `a`, pull one value from `b` for each, and yield the two together, returning when `b` runs out or `yield` returns false, with `b` stopped on every way out.
