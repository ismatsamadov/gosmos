---
slug: hashing-a-stream
title: Hashing a Stream
type: code_tests
position: 2
difficulty: 6
topics: crypto, streams
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Tuesday. The program that delivers files, dlinkd, gets its first fingerprints today: one for each Furrow-A image, 10.9 GB apiece, without ever holding a whole image in memory.

`sha256.New` returns a `hash.Hash`: a writer that digests whatever you write into it, piece by piece, and whose `Sum(nil)` returns the digest of everything written so far. Written in pieces, a text gets the same digest as written at once:

```go
h := sha256.New()
h.Write([]byte("Apogee ground "))
h.Write([]byte("calling Pipit-1"))
fmt.Printf("%x\n", h.Sum(nil)[:4]) // fe5e0f73
```

## Your task

`fingerprintStream` hashes only what one `Read` returns, the first 4 KiB of a long stream. Make it `io.Copy` the whole stream into `sha256.New()` and return the hex of `Sum(nil)`, or `""` and the copy's error.

> [!NOTE]
> A stream that breaks off returns an error from `io.Copy`. A fingerprint of part of a file is worse than none, so return no digest at all.
