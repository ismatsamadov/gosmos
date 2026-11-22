---
slug: digests-that-name-their-hash
title: Digests That Name Their Hash
type: code_tests
position: 4
difficulty: 7
topics: crypto, errors, strings
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Cirrine's archive fingerprints with SHA3-256, Apogee's deliveries with SHA-256, and both write 64 hex digits, so a bare digest could be either. Mirela's rule for the folder where deliveries land: every digest names the hash that made it, as in `sha256:fe5e…bd66`.

Read such text strictly, so nothing is ever guessed: a name you know, spelled exactly, then exactly 64 lowercase hex digits. `strings.Cut` splits at the first `:` and says whether there was one:

```go
a, h, ok := strings.Cut("crc32:f6ced6c0", ":")
fmt.Println(a, h, ok) // crc32 f6ced6c0 true
```

## Your task

`parseDigest` splits at the first `:` and decodes whatever follows. Make it refuse, in this order, text with no `:`, any name but `sha256` and `sha3-256`, and anything after the colon but 64 lowercase hex digits, with these errors:

```text
digest must be <hash>:<hex>
unknown hash "md5"
digest must be 64 lowercase hex digits
```
