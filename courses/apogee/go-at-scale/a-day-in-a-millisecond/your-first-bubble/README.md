---
slug: your-first-bubble
title: Your First Bubble
type: code_tests
position: 1
difficulty: 9
topics: context, testing, time
created: 2026-11-11T16:51:57Z
modified: 2026-11-12T08:43:36Z
---
Tuesday 24 April. Nkem: "You'll test timers for the rest of the season. Write one bubble yourself."

`synctest.Test(t, f)` runs `f` in a bubble, whose goroutines share a fake clock. It starts at midnight UTC on 1 January 2000, and whenever all of them are blocked, on a timer, a `Sleep` or a channel, it jumps to the next timer due: a four-hour wait ends at once.

A test shifts it to take-over day:

```go
synctest.Test(t, func(t *testing.T) {
	shift := day.Sub(time.Now()) // day: 28 June
	now := func() time.Time {
		return time.Now().Add(shift).UTC()
	}
	time.Sleep(4 * time.Hour) // no real time
	fmt.Println(now())        // 04:00 on 28 June
})
```

## Your task

`checkWaitUntil`, in `checks.go`, checks `wait` on the real clock. Run the checks in a bubble on 28 June's clock, and fail `t` unless a wait for 04:12:25 returns `nil` at 04:12:25 exactly, a wait for 03:12:25 then returns `nil` at once, and a wait for 05:12:25 with a context already canceled returns `context.Canceled`.

> [!NOTE]
> The tests run your check on your `waitUntil`, then on four mutants, each with one bug: every mutant must fail it.
