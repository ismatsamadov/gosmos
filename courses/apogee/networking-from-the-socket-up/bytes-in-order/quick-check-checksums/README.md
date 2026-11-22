---
slug: quick-check-checksums
title: Quick Check: Checksums
type: quiz
position: 7
difficulty: 7
topics: bytes, errors
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Thursday 30 March. Mirela reads the new frame code from Port Ellery. "A CRC tells you about noise. It tells you nothing about people." Then three questions about this check:

```go
n := len(payload) - 4
want := binary.BigEndian.Uint32(payload[n:])
if crc32.ChecksumIEEE(payload[:n]) != want {
	return Frame{}, ErrBadCRC
}
```

## Your task

Answer Mirela's three questions: what the check does with a deliberate change, why it comes before any field is read, and whether it is the CRC-32 that downlinked files carry.
