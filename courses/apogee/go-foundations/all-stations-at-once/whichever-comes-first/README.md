---
slug: whichever-comes-first
title: Whichever Comes First
type: code_tests
position: 6
difficulty: 7
topics: channels, select
created: 2026-11-01T07:36:47Z
modified: 2026-11-01T07:36:47Z
---
Tonight Thistle may rise over Lantern Atoll or over Cape Wren first, and Professor Marr wants its first beacon, from whichever station hears it.

`select` waits on several channel operations at once and runs the case of whichever is ready first. If more than one is ready, it picks one of them at random.

```go
select {
case kmh := <-wind:
	fmt.Println("wind first:", kmh)
case mm := <-rain:
	fmt.Println("rain first:", mm)
}
```

## Your task

`first` waits only on `lan`, so a beacon heard at Cape Wren is never seen. Make it wait on both channels with `select`, returning whichever beacon comes first.
