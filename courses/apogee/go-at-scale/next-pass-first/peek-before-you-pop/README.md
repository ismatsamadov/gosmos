---
slug: peek-before-you-pop
title: Peek Before You Pop
type: code_tests
position: 5
difficulty: 8
topics: generics, methods
created: 2026-11-11T12:51:43Z
modified: 2026-11-11T12:51:43Z
---
Ingrid's corner of the test Board shows each station's next pass and the minutes until its AOS, and that pass must stay in the queue until a dish has flown it.

`container/heap` keeps the smallest value at index 0 of its slice, so reading it costs nothing and moves nothing. An empty queue has no first value; then `Peek` returns the zero value of `T`, which `var zero T` declares whatever `T` is: 0 for an int, "" for a string.

```go
h := &gusts{72, 50, 81}
heap.Init(h)
fmt.Println((*h)[0], h.Len()) // 50 3
```

## Your task

`Peek` is a stub. Return the queue's smallest value and true, leaving it in the queue, or the zero value and false when the queue is empty.

> [!TIP]
> The heap's slice is `q.h.items`.
