---
slug: only-the-first-part
title: Only the First Part
type: code_tests
position: 6
difficulty: 3
topics: interfaces, streams, testing
created: 2026-11-04T17:36:06Z
modified: 2026-11-04T17:36:06Z
---
The recorder's quick preview must stop after a small prefix. Paz still needs the remaining bytes for the full report, so reading everything and slicing a string afterwards would consume too much.

`io.LimitReader(src, n)` supplies at most `n` bytes from its source. Its caller sees a normal end when the limit is reached. Zero or a negative limit supplies no bytes:

```go
prefix := io.LimitReader(source, 16)
data, err := io.ReadAll(prefix)
```

## Your task

Implement `preview` by reading a limited reader with `io.ReadAll`. Return text and the read error. Do not close or drain `src`: the caller will continue reading it after the preview.
