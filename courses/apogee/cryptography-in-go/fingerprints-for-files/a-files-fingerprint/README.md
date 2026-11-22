---
slug: a-files-fingerprint
title: A File's Fingerprint
type: code
position: 0
difficulty: 6
topics: crypto, formatting
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T07:46:46Z
---
Monday 17 July. Before Yusra promises Cirrine a fingerprint for every file, Tomás tries one on the first words you ever sent up, `Apogee ground calling Pipit-1`, and on the same words for Pipit-2: "Twenty-nine bytes in. Let's see what comes out!"

A recipe that turns any bytes into a fingerprint of fixed size is a hash, and what it returns is a digest. The hash Apogee uses, SHA-256, makes a 32-byte digest, whether the input is 29 bytes or a whole image, and one changed byte gives a digest that looks unrelated. `sha256.Sum256` returns a `[32]byte`:

```go
sum := sha256.Sum256([]byte("PING"))
fmt.Printf("%x\n", sum[:4]) // 906055e5
```

## Your task

`main` prints each text's CRC-32, the 4-byte check number every file from a satellite carries. Make it print `sha256.Size` as the digest's size, then the first 8 bytes of each text's `sha256.Sum256` in hex:

```text
digest size: 32 bytes
Pipit-1: fe5e0f73525ffc9f
Pipit-2: 8cbc29efabfcd2fb
```
