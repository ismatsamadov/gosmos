---
slug: the-length-comes-first
title: The Length Comes First
type: code_tests
position: 2
difficulty: 7
topics: bytes, networking, streams
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
After its banner, a controller speaks only in frames: messages of raw bytes, each starting with its length. Nkem shows why on the WELCOME that answers dishd's HELLO: when the last command was number 10, a `0x0a`, which is a newline, sits in its middle.

The first two bytes of a frame are its length, high byte first, and `binary.BigEndian.Uint16` turns them into a number: how many bytes follow. Read the two, then exactly that many:

```go
b := []byte{0x00, 0x0e}
n := binary.BigEndian.Uint16(b)
fmt.Println(n) // 14
```

## Your task

`readPayload` returns `nil`. Read the two length bytes with `io.ReadFull`, turn them into the length, then read exactly that many bytes and return them, or the error that stopped you.

> [!WARNING]
> Never look for a frame's end with `ReadBytes('\n')`: it stops at the first `0x0a`, wherever that is.
