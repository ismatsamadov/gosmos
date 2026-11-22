---
slug: refuse-a-damaged-frame
title: Refuse a Damaged Frame
type: code_tests
position: 6
difficulty: 8
topics: bytes, errors
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Wednesday 29 March. Sunniva's cable test flipped one bit of a STATUS on its way in, and dishd showed Dish 1 at 203.73° while it sat at 183.25°.

A CRC is worth something only if it is checked before the frame is used. Compute the CRC-32 of everything but the last four bytes and compare it with those four, read big-endian. If they differ, some bit changed on the way and no field can be trusted, so return an empty `Frame` with the error:

```go
want := binary.BigEndian.Uint32(trailer)
if crc32.ChecksumIEEE(file) != want {
	return fmt.Errorf("%s: damaged", name)
}
```

## Your task

`decodeFrame` splits the payload without checking it. Refuse a payload under 9 bytes with `ErrBadLength` wrapped with its length, and one whose last four bytes don't match the CRC-32 of the rest with `ErrBadCRC`. Return an empty `Frame` with either.

> [!NOTE]
> A CRC-32 catches every single flipped bit, wherever it lands, the CRC's own bits included.
