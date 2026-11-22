---
slug: quick-check-where-garbage-comes
title: Quick Check: Where Garbage Comes From
type: quiz
position: 2
difficulty: 8
topics: conversions, performance, pointers
created: 2026-11-12T03:10:33Z
modified: 2026-11-12T03:10:33Z
---
Tuesday 22 May. Nkem writes four lines on the whiteboard and asks which of them make garbage.

A value goes on the heap, and later becomes garbage, when it can't stay in its function's stack frame: when it may outlive the call, say, or its size is known only at run time. In a terminal, `go build -gcflags=-m` prints what moved there, such as `moved to heap: c`.

```go
func note(c Contact, b []byte) (*Contact, int) {
	p := &c                   // 1
	dish := c.Dish            // 2
	fmt.Println(c.MaxElevDeg) // 3
	p.ID = string(b)          // 4
	return p, dish
}
```

## Your task

Answer Nkem's three questions about the four lines.
