---
slug: coarse-steps-then-seconds
title: Coarse Steps, Then Seconds
type: code_tests
position: 3
difficulty: 9
topics: profiling, time
created: 2026-11-10T15:36:39Z
modified: 2026-11-11T17:33:39Z
---
Thursday. 81% of passplan's time goes to `elevationDeg`, asked about every second for every satellite and station. Tomás: "A satellite takes minutes to climb into view! Ask every 30 seconds, then look closer."

The code a program spends most of its time in is its hot path, and the first fix is to run it less. Where two answers differ, `crossing`, in `passplan.go`, halves the gap between them down to a millionth of a second, and a wide gap costs little more: 30 seconds take 25 halvings, where one second takes 20.

```go
gap, n := 30*time.Second, 0
for gap > time.Microsecond {
	gap /= 2
	n++
}
fmt.Println(n, gap)
// 25 894ns
```

## Your task

Make `findPasses` ask `el` every 30 seconds instead of every second, then once more at the last second before `to`, where the old scan stopped. Give `crossing` the two moments whose answers differ. The test wants the old scan's passes from at most a twentieth of its questions.

> [!NOTE]
> The plan keeps only passes of 2 minutes or more, and none of those fits between two questions 30 seconds apart.
