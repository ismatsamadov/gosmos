---
slug: number-every-command
title: Number Every Command
type: code_tests
position: 4
difficulty: 8
topics: methods, networking
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Thursday 11 May. Tomás asks from Port Ellery the question Nkem asked about Pipit-1 a year ago: "What does the dish see if this line runs twice?"

A dish controller runs a command only if its sequence number is above `last_seq`, the number of the last command it ran, kept through power cuts; it refuses any other as a duplicate. So dishd numbers commands upward from `l.next`, and after each HELLO goes on from one past the WELCOME's `LastSeq`, never back:

```go
peak := 0.0
for _, el := range []float64{12.5, 61.4, 33.0} {
	peak = max(peak, el)
}
fmt.Println(peak) // 61.4
```

## Your task

In `link.go`, make `number` set the frame's `Seq` to `l.next` and count `l.next` up, and make `learn` set `l.next` to the larger of itself and `w.LastSeq + 1`.

> [!WARNING]
> `l.next = w.LastSeq + 1` alone goes back when a controller's count is behind dishd's, to numbers dishd has already used.
