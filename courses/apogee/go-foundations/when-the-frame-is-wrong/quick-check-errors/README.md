---
slug: quick-check-errors
title: Quick Check: Errors
type: quiz
position: 3
difficulty: 4
topics: errors
created: 2026-11-01T07:36:46Z
modified: 2026-11-22T10:08:26Z
---
The new log goes live tonight. Before it does, Ingrid runs her quick check: three questions, a few seconds each. The second one is about this code from the ingest loop, where `store` saves a reading for the customer:

```go
err := checkBattery("pipit-1", 7.9)
if err != nil {
	fmt.Println("rejected:", err)
	continue
}
store(7.9)
```

## Your task

Pick one answer for each question, then press **Submit**.
