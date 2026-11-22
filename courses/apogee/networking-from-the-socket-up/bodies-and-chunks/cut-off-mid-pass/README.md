---
slug: cut-off-mid-pass
title: Cut Off Mid-Pass
type: code_tests
position: 4
difficulty: 9
topics: errors, http, streams
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T19:06:41Z
---
Tuesday 6 June, 08:06. The recorder restarted in the middle of Furrow-A's 08:02 pass, and its stream simply stopped. Sunniva: "Wind light. Recorder restarted. Half a pass kept as whole."

A chunked body has exactly one end: the zero chunk and the blank line after its trailer. A stream that ends anywhere before that is a cut recording, however much it delivered. A read reports the end as `io.EOF`, or as `io.ErrUnexpectedEOF` when it ends partway through what was asked for:

```go
br := bufio.NewReader(strings.NewReader("10"))
_, err := readLine(br)
fmt.Println(err) // unexpected EOF
```

## Your task

`readChunked` takes the end of the stream between two chunks for the end of the body. Answer 400 `recording cut off` instead for an end anywhere before the final blank line: in a size line, in a chunk, after the zero chunk or in the trailer. The given `put` then keeps nothing.
