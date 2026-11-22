---
slug: what-the-tag-adds
title: What the Tag Adds
type: code
position: 1
difficulty: 7
topics: crypto, printing
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Yusra is writing to Cirrine about sealed commands and wants numbers, not adjectives: "TL;DR: how many bytes longer, and what happens to a changed one?"

A GCM value is a `cipher.AEAD`. Its `NonceSize()` is the size of the nonce every seal needs, and its `Overhead()` is what `Seal` adds: the tag. When `Open` refuses, its error says only that the check failed. To change sealed bytes and keep the original, change a copy:

```go
line := []byte("PING")
bad := bytes.Clone(line) // a copy to change
bad[len(bad)-1] ^= 1     // its last byte
fmt.Printf("%s %s\n", line, bad) // PING PINF
```

## Your task

`main` seals `SET_MODE|SAFE` with Cirrine-1's test seal key and the tests' fixed nonce, then prints zeros. Print both lengths, `NonceSize()` and `Overhead()`, and the error from opening a copy of the sealed bytes with its last byte changed:

```text
plaintext: 13 bytes
sealed: 29 bytes
nonce: 12 bytes, tag: 16 bytes
tampered: cipher: message authentication failed
```

> [!NOTE]
> A fixed nonce is for test values only. Every real seal takes a fresh one, as Wednesday shows.
