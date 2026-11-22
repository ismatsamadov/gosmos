---
slug: the-crc-comes-last
title: The CRC Comes Last
type: code_tests
position: 7
difficulty: 9
topics: errors, http, strings
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T19:06:41Z
---
Thursday afternoon. Nkem notes that the recorder knows a recording's CRC, the checksum that shows whether a byte changed on the way, only at LOS, when the satellite sets. So it sends the CRC last, the way a satellite sends a file's CRC after the file.

A chunked upload can end with trailer fields after its zero chunk: Go's client sends `Apogee-Crc32: e2fc055d` there when the request's `Trailer` is set. An upload with a `Content-Length` can carry the same field in its header instead. The value is 8 hex digits, in either case:

```go
n, err := strconv.ParseUint("A2912082", 16, 32)
fmt.Printf("%08x %v\n", n, err) // a2912082 <nil>
```

## Your task

`checkCRC` returns the body's CRC without checking it. Take `apogee-crc32` from the trailer, or else from the header; with neither, accept the body. Refuse a value that isn't 8 hex digits, and one that doesn't match, with 400:

```text
bad apogee-crc32 "xyz"
crc mismatch: sent a2912083, received a2912082
```
