---
slug: write-to-any-destination
title: Write to Any Destination
type: code_tests
position: 1
difficulty: 3
topics: interfaces, streams, testing
created: 2026-11-04T17:36:06Z
modified: 2026-11-07T19:09:32Z
---
Paz wants the van's status report on the screen during a check and in a file afterwards. A writer that always reaches for a filename cannot do both without another copy of the code.

`io.Writer` describes something that accepts bytes. `fmt.Fprintf` writes formatted text to the writer supplied first and returns a byte count and error:

```go
var report strings.Builder
_, err := fmt.Fprintf(&report, "passes: %d\n", 3)
fmt.Println(err)
```

## Your task

Implement `writeStatus`. Write exactly `code: state` followed by one newline to `w`, with the two arguments substituted as text. Return the writing error; ignore only the byte count. The tests use an in-memory writer and one that refuses writes.
