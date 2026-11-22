---
slug: a-struct-like-the-frame
title: A Struct Like the Frame
type: code_tests
position: 3
difficulty: 7
topics: bytes, errors, structs
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Sunniva wants every field of the status frame on the hut's screen, not only the azimuth and the temperature: the dish's state, the wind, the flags and the number of the last command it carried out.

`binary.Decode(buf, order, &v)` fills a struct field by field, in order, as long as every field has a fixed size, such as `uint8`, `int16` or `uint32`. It reads as many bytes as the struct needs and ignores any after them, so check the length yourself first:

```go
type reading struct {
	Station uint8
	TempC10 int16
}
var r reading
n, err := binary.Decode([]byte{3, 0xff, 0x83},
	binary.BigEndian, &r)
fmt.Println(n, err, r) // 3 <nil> {3 -125}
```

## Your task

`decodeStatus` returns an empty `Status`. Refuse a body that isn't 16 bytes with an error giving its length:

```text
status body is 15 bytes, want 16
```

Then decode it into the given `rawStatus`, big-endian, and fill the `Status` in degrees, km/h and °C.
