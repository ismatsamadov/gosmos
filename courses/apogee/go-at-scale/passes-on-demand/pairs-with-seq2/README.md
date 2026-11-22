---
slug: pairs-with-seq2
title: Pairs with Seq2
type: code_tests
position: 3
difficulty: 8
topics: iterators, methods
created: 2026-11-11T21:35:30Z
modified: 2026-11-12T08:43:35Z
---
Wednesday 28 March. Tirrin Rock (TIR) goes live on Monday with two dishes, and Ingrid wants its panel on the Board to list each contact beside the dish that flies it.

An `iter.Seq2[K, V]` yields two values at a time, and a range loop with two variables takes both. Inside, the iterator calls `yield` with both and stops when that returns false, as a one-value iterator does. `slices.All` is one, yielding each index with its value:

```go
codes := []string{"TIR", "VSK", "CAL"}
for i, code := range slices.All(codes) {
	fmt.Println(i, code) // 0 TIR, 1 VSK, 2 CAL
}
```

A method can return one too, so the panel ranges over `sch.At("TIR")` and gets each dish with its contact.

## Your task

`At` yields nothing. Make it yield the dish and the contact of each of the station's contacts, in the schedule's order.

> [!NOTE]
> The tests use Monday's schedule as passplan 2 plans it: Tirrin Rock has 48 contacts, 13 of them on dish 2.
