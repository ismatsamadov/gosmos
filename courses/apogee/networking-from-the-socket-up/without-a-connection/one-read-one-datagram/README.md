---
slug: one-read-one-datagram
title: One Read, One Datagram
type: code_tests
position: 1
difficulty: 8
topics: errors, networking, slices
created: 2026-11-13T06:15:45Z
modified: 2026-11-13T16:03:23Z
---
Monday, 07:38, Furrow-A's pass: every line on Ingrid's screen in the Fishbowl ended at `mode`. Ingrid: "Negative. Frames cut short." The live view's read had room for 64 bytes, and a frame is 72.

On a TCP stream, bytes a read leaves behind wait for the next read. A datagram never waits: each `ReadFrom` returns exactly one, and when your slice is too short it gets the datagram's start, the rest is thrown away, and there is no error. So read into a slice one byte longer than the biggest datagram you accept: a datagram that fills it has been cut.

```go
buf := make([]byte, 4)
pc.WriteTo([]byte("PING"), pc.LocalAddr())
n, _, err := pc.ReadFrom(buf)
fmt.Println(n, err) // 4 <nil>: whole, or cut?
```

## Your task

`readDatagram` in `datagram.go` reads into 64 bytes. Read into `maxDatagram+1` instead, and when more than `maxDatagram` (512) bytes arrive, return `nil`, the sender and `ErrTooBig`.

> [!NOTE]
> A datagram can be empty: `ReadFrom` then returns 0 bytes and no error, and so does `readDatagram`.
