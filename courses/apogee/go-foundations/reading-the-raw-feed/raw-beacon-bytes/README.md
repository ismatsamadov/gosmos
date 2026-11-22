---
slug: raw-beacon-bytes
title: Raw Beacon Bytes
type: code_tests
position: 8
difficulty: 6
topics: bytes, errors
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:27Z
---
Pipit beacons reach the Fishbowl as 24 hex digits: 12 raw bytes. Ingrid wants each beacon's sequence number and battery on the Board.

`hex.DecodeString` turns hex text into a `[]byte`, or returns an error if the text isn't hex. A number stored in two bytes, high byte first, is big-endian, and `binary.BigEndian.Uint16` reads one from a two-byte slice:

```go
b, err := hex.DecodeString("0bb8")
fmt.Println(b, err) // [11 184] <nil>
n := binary.BigEndian.Uint16(b)
fmt.Println(n) // 3000
```

## Your task

Write `decodeBeacon`. Return an error if the text isn't hex, isn't exactly 12 bytes, or doesn't start with `AP`, the sync bytes that mark a beacon. Otherwise return the sequence number from bytes 3 and 4, and the battery millivolts from bytes 6 and 7.

> [!TIP]
> `b[3:5]` holds bytes 3 and 4: a slice stops before its second index.
