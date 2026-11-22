---
slug: keep-the-last-error
title: Keep the Last Error
type: code_tests
position: 1
difficulty: 4
topics: errors, streams, testing
created: 2026-11-06T16:34:04Z
modified: 2026-11-06T16:34:04Z
---
Paz's backup of the log sits on her oldest memory stick, and reading it fails partway. She wants the lines that arrived and the reason the rest did not, never a short list that looks complete.

`Scan` returns false at the end of the input and also when reading fails, so the loop alone cannot tell which happened. `sc.Err()` can: it returns nil after a clean end, never `io.EOF`, and the failure otherwise. The lines scanned before a failure did arrive, though the last one may be cut short.

```go
for sc.Scan() {
    frames++
}
if err := sc.Err(); err != nil {
    fmt.Println("stopped after", frames, err)
}
```

## Your task

Replace the TODO in `collectLines`: return the lines together with the scanner's error, so a failure brings back the lines read so far and its reason.
