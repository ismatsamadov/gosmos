---
slug: ten-minutes-of-predictions
title: Ten Minutes of Predictions
type: code
position: 2
difficulty: 3
topics: loops
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:19:30Z
---
Tomás has computed Pipit-3's angle above Ystmark for each minute of the coming pass. This forecast, a prediction, is ready to check. Tomás wants the first ten minutes on the console, from `elevationAt` in `prediction.go`.

When you know how many rounds you need, `for i := range n` is the shortest loop: it runs n times, with `i` counting 0, 1, 2 and so on, up to `n-1`.

```go
for dish := range 2 {
	fmt.Println("Ystmark dish", dish+1, "parked")
}
```

## Your task

In `main`, replace `minute := 0` with a `range` loop around the `Printf`, so it prints minutes 0 to 9.
