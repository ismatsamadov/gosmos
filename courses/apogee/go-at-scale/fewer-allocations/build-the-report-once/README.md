---
slug: build-the-report-once
title: Build the Report Once
type: code_tests
position: 1
difficulty: 8
topics: performance, strings
created: 2026-11-12T03:10:33Z
modified: 2026-11-12T09:50:21Z
---
Monday afternoon. On every replan each host gets a text report of its station's contacts, and Cape Wren's 76 lines are joined with `+=`: every line copies the whole report so far, 78 KB of copying for 2 KB of text.

A string never changes, so `+=` makes a new one and copies the old text and the new line into it. A `strings.Builder` writes into one buffer instead, and `String` hands that buffer over without a copy. `Grow` the builder by the whole length first, and the buffer is allocated once:

```go
var b strings.Builder
b.Grow(len(codes) * 4)
for _, code := range codes {
	b.WriteString(code)
	b.WriteByte(' ')
}
fmt.Println(b.String()) // YST PEL SQT ...
```

## Your task

Build `stationReport`'s text in a `strings.Builder`. Add up the length of the heading and of every line first, then `Grow` the builder by that: Cape Wren's report should take at most 2 allocations instead of 77.

> [!TIP]
> Every line after the heading is the pass ID, a space, the word `dish`, a space and the dish's number, 1 or 2, then a newline.
