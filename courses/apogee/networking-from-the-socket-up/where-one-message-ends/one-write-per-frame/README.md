---
slug: one-write-per-frame
title: One Write per Frame
type: code_tests
position: 6
difficulty: 7
topics: goroutines, networking, streams
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Thursday night. Two of dishd's goroutines sent frames at the same moment, and the controller read the first half of one glued to the second half of the other.

A connection writes each `Write` call whole, even when goroutines share it, but nothing keeps two calls together: another goroutine's frame can land between them. Build the whole frame in memory and send it with one `Write`. Then return its error, which is the only sign that the line has failed:

```go
var b bytes.Buffer
b.WriteString("YST-1 ")
b.WriteString("TRACKING\n")
_, err := w.Write(b.Bytes()) // one Write
return err
```

## Your task

`writeFrame` writes the length, then the rest, in two `Write`s, and ignores both errors. Send the given `encodeFrame(f)` in a single `Write` and return its error.
