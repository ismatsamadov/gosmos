---
slug: the-sealed-line
title: The Sealed Line
type: code_tests
position: 5
difficulty: 8
topics: crypto, errors, strings
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Tuesday. Sunniva, from Ystmark: "Wind 41. Queue takes lines. One command, one line." A sealed command must cross to Ystmark's uplinkd as one line of text and read back there exactly as it left.

The line has seven fields with single spaces between them: `SEAL`, the version `2`, the satellite, the number, the key ID, then the nonce and the sealed bytes in lowercase hex. `strings.Split` on one space keeps the empty field a doubled space makes, where `strings.Fields` would hide it, and `%x` writes bytes as hex:

```go
f := strings.Split("PING  YST", " ")
fmt.Println(len(f))              // 3, one empty
fmt.Printf("%x\n", []byte("Hi")) // 4869
```

## Your task

1. Make `formatSealed` write `s` as its line.
2. Make `parseSealed` read a line back. Wrap `ErrVersion` for a `SEAL` line of any version but `2`, and `ErrMalformed` for any other fault: not seven fields, no `SEAL`, an ID `validID` refuses, a number below 1, a nonce that isn't 12 bytes or a ciphertext that isn't 17 to 272, as `hexBytes` checks them.

> [!NOTE]
> 17 to 272 bytes is a command of 1 to 256 bytes, plus its 16-byte tag.
