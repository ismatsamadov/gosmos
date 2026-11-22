---
slug: all-eleven-bytes
title: All Eleven Bytes
type: code_tests
position: 1
difficulty: 7
topics: networking, streams
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Tuesday. Over Sunniva's test line, which runs from the simulator to Port Ellery and back over the backhaul, the 11-byte PONG a controller sends when dishd asks "are you there?" with a PING arrives in pieces, and the dish service's single `Read` takes the first three bytes for the whole answer.

A connection is a stream of bytes, and one `Read` returns whatever has arrived: all of an answer, or three bytes of it. `io.ReadFull(r, buf)` reads until `buf` is full. If the stream ends first, it returns `io.ErrUnexpectedEOF`, or `io.EOF` when nothing came at all:

```go
r := iotest.OneByteReader(
	strings.NewReader("CRATE"))
buf := make([]byte, 5)
n, _ := r.Read(buf)
fmt.Println(n) // 1
n, _ = io.ReadFull(r, buf[n:])
fmt.Println(n) // 4
```

## Your task

`readExactly` makes one `Read`. Make it read exactly `n` bytes with `io.ReadFull`, returning the error when the stream ends too soon.
