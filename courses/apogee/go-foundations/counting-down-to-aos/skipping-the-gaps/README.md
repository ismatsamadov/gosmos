---
slug: skipping-the-gaps
title: Skipping the Gaps
type: code_tests
position: 5
difficulty: 4
topics: loops
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:07:33Z
---
Pipit-3's first battery readings are in, with holes: some frames, the individual pieces of data sent over the radio, were lost in noise. Lost frames read as `lostFrame`, or -1. Average the holes in, and a healthy battery looks sick.

`continue` skips the rest of the current round and moves on to the next one:

```go
for station := range 5 {
	if station == 4 {
		continue // Cape Wren is stowed
	}
	fmt.Println("polling station", station)
}
```

## Your task

Make `averageBattery` average the readings from minute `from` up to, but not including, `to`. Skip lost frames with `continue`, and return 0 when no reading is left.

> [!TIP]
> In a three-part `for`, the update, such as `i++`, still runs after `continue`. A counter you raise at the end of the loop's block is skipped instead, and the loop can run forever.
