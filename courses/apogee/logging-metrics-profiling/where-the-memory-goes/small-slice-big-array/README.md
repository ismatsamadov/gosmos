---
slug: small-slice-big-array
title: Small Slice, Big Array
type: code_tests
position: 4
difficulty: 9
topics: profiling, slices
created: 2026-11-10T15:36:39Z
modified: 2026-11-11T17:33:39Z
---
Wednesday, 11:00. Reading on from `runPass`, Nkem stops at the pass report: it keeps the last 64 bytes the radio heard, sliced straight from the pass's 4 MiB buffer. "Sixty-four bytes. How much memory do they keep?"

A slice is a window onto an array, and while any window onto it is kept, the whole array stays in use: 64 bytes sliced from a 4 MiB buffer keep all 4 MiB. Copying them into a slice of their own lets the buffer go:

```go
big := make([]byte, 1<<20)
small := big[:8]          // a window onto 1 MiB
own := bytes.Clone(small) // 8 bytes of its own
fmt.Println(cap(small), cap(own))
// 1048576 8
```

## Your task

`tail` returns a window onto `frames`. Make it return a copy of the last `n` bytes, or of all of `frames` when there are fewer than `n`.

> [!NOTE]
> A copy also stays as it was when the buffer is written to later, which a window onto the buffer doesn't.
