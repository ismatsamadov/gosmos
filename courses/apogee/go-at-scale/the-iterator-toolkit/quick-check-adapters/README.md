---
slug: quick-check-adapters
title: Quick Check: Adapters
type: quiz
position: 3
difficulty: 8
topics: generics, iterators
created: 2026-11-11T22:52:57Z
modified: 2026-11-11T22:52:57Z
---
Before Ingrid's corner of the Board goes live, Nkem writes its chain on the whiteboard and asks what it does, and when.

`contacts` is Monday's schedule as a sequence, 260 contacts at ten stations; `atTIR` keeps Tirrin Rock's, and `line` turns a contact into a line of text:

```go
lines := Take(Map(Filter(contacts, atTIR),
	line), 3)
for l := range lines {
	fmt.Println(l)
}
```

## Your task

Answer Nkem's three questions: what has run once `lines` is built, in what order the work happens, and what stops once `Take` has its third line.
