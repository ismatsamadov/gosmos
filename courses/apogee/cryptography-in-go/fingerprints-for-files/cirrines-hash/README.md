---
slug: cirrines-hash
title: Cirrine's Hash
type: code_tests
position: 5
difficulty: 7
topics: crypto, interfaces
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Cirrine's archive has checked files with SHA3-256 since its first satellite, and Yusra's message is one line: "TL;DR: Cirrine checks our files with their hash, not ours."

SHA3-256 is a different design from SHA-256, with a digest of the same size, 32 bytes. `sha3.New256` returns it as a `hash.Hash`, the interface every hash in Go's library shares, so code that writes and sums through `hash.Hash` works with either:

```go
var h hash.Hash = sha3.New256()
io.WriteString(h, "PING")
fmt.Println(h.Size()) // 32
```

## Your task

`newHash` knows only `sha256`. Make it return `sha3.New256()` for `sha3-256`, so that `fingerprintAs`, in the read-only `digest.go`, can name and compute either digest. Any other name must still fail before a single byte is read.
