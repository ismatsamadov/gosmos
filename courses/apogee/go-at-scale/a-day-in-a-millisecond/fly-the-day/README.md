---
slug: fly-the-day
title: Fly the Day
type: code_tests
position: 7
difficulty: 9
topics: context, testing, time
created: 2026-11-11T16:51:57Z
modified: 2026-11-11T16:51:57Z
---
Friday 4 May. The second wave of stations goes live on Monday, and Ingrid wants the station runner flown through take-over day's plan, in a test, before any host's dish follows it.

The last two weeks' pieces fit together: a loop over the contacts, a wait for each AOS, and a context that can stop it. In a bubble, a whole day of contacts flies in milliseconds of real time, and a cancel set for noon lands at 12:00:00 exactly:

```go
synctest.Test(t, func(t *testing.T) {
	now := clockAt(takeOverDay)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	time.AfterFunc(12*time.Hour, cancel) // noon
	err := flyDay(ctx, now, day, fly)
	fmt.Println(err, now()) // canceled at 12:00
})
```

## Your task

`flyDay` calls `fly` for every contact at once. For each contact, in the order given, wait for its AOS with `waitUntil` before you fly it, and return the wait's error as soon as one fails.

> [!NOTE]
> The tests plan 28 June from scratch, with every dish free at midnight: 2,020 contacts, 992 of them before noon.
