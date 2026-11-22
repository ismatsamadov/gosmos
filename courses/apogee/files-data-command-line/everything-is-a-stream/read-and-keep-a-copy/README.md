---
slug: read-and-keep-a-copy
title: Read and Keep a Copy
type: code_tests
position: 5
difficulty: 3
topics: interfaces, streams, testing
created: 2026-11-04T17:36:06Z
modified: 2026-11-12T11:41:59Z
---
Paz needs to see exactly what a decoder read when a recording stops halfway through. Copying a file afterwards would not show which bytes the decoder actually consumed.

`io.TeeReader(src, log)` is a reader that writes each consumed block to its second argument before returning it. It does no reading by itself and keeps no extra buffer. A log-writing failure becomes a read error for its caller:

```go
traced := io.TeeReader(source, &trace)
data, err := io.ReadAll(traced)
```

## Your task

Implement `captureText` by reading from a `TeeReader` with `io.ReadAll`. Return the consumed bytes as text and preserve the error. A source that fails after a partial block should leave those consumed bytes in the trace as well.
