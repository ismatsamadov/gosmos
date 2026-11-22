---
slug: range-over-a-function
title: Range Over a Function
type: code
position: 0
difficulty: 8
topics: functions, iterators, time
created: 2026-11-11T21:35:30Z
modified: 2026-11-12T08:43:35Z
---
Hollis, who runs operations at Furrowlight, wants every day of his summer listed before he books a pass in one of them.

Go can range over a function. An iterator, `iter.Seq[V]`, is a function that takes another, `yield`, and calls it with each value in turn: each call runs the loop's body once. When the loop stops early, `yield` returns false, and the iterator must return:

```go
func threeTwoOne(yield func(int) bool) {
	for n := 3; n > 0; n-- {
		if !yield(n) {
			return
		}
	}
}
```

A loop `for n := range threeTwoOne` gets 3, 2, then 1.

## Your task

`Days` returns an iterator that yields nothing. Make it yield `n` midnights, a day apart, starting at `from`. `main` prints the first week of the third quarter, July to September:

```text
Sun 1 Jul
Mon 2 Jul
Tue 3 Jul
Wed 4 Jul
Thu 5 Jul
Fri 6 Jul
Sat 7 Jul
```

> [!TIP]
> `from.AddDate(0, 0, i)` is the midnight `i` days after `from`.
