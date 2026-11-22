---
slug: noise-is-not-an-attacker
title: Noise Is Not an Attacker
type: code_tests
position: 1
difficulty: 6
topics: crypto, strings
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Mirela opens a Furrowlight field report, the line that says how green field 12 looks from orbit, changes 0.71 to 0.17 and adds four bytes. The CRC-32 its delivery note carries still matches: "A CRC catches what noise does by chance. I don't change things by chance."

A CRC is built to catch radio noise, and its arithmetic runs backward: four chosen bytes bring any file back to any CRC-32, as Mirela's do:

```go
a := []byte("NDVI field 12: 0.71\n")
b := append([]byte("NDVI field 12: 0.17\n"),
	0x1a, 0x15, 0xbc, 0xd1)
fmt.Printf("%08x %08x\n",
	crc32.ChecksumIEEE(a), crc32.ChecksumIEEE(b))
// f6ced6c0 f6ced6c0
```

## Your task

`fingerprint` returns the file's CRC-32 as 8 hex digits. Make it return the 64 hex digits of its `sha256.Sum256`: nobody knows a way to make a second file with the same SHA-256 digest. The test plays Mirela, whose forgery must get a fingerprint of its own.
