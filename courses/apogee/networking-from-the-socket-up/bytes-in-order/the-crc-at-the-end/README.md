---
slug: the-crc-at-the-end
title: The CRC at the End
type: code_tests
position: 5
difficulty: 8
topics: bytes, networking
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Tuesday afternoon. Tomás sends dishd's first real POINT to the simulator, and it answers with a refusal, a NAK, reading `bad crc`: the frame's last four bytes, its CRC, a check number computed from the rest, are all zero.

`crc32.ChecksumIEEE(b)` computes the CRC-32 of exactly the bytes you give it, and the protocol names them: the type, seq and body, everything after the length. Not the length, and not the CRC itself. The controller computes it again and refuses a frame whose numbers disagree:

```go
sum := crc32.ChecksumIEEE([]byte("YST-1"))
b := binary.BigEndian.AppendUint32(nil, sum)
fmt.Printf("% x\n", b) // 30 d8 26 52
```

## Your task

`encodeFrame` ends every frame with a CRC of zero. Replace the `TODO` so it appends the CRC-32 of the type, seq and body, the bytes after the length, big-endian.

> [!WARNING]
> Leave the length out: a PING whose CRC covers it ends `c0 a3 dd d6`, and the controller refuses it.
