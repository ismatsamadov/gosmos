---
slug: proof-for-furrowlight
title: Proof for Furrowlight
type: code_tests
position: 7
difficulty: 9
topics: crypto, streams, time
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Friday 1 September. Furrowlight checks the receipt for Sunday's Furrow-A file against the file already in its own storage, with the checker from Yusra's guide, and Hollis is waiting.

Each check believes only what the checks before it have proved. A key signs only inside its window: from `Activate`, up to but not including `Retire`, and a zero `Retire` means it still signs. `io.Copy` returns the number of bytes it copied, so hashing a file also measures it:

```go
h := sha256.New()
n, _ := io.Copy(h, strings.NewReader("PING"))
fmt.Println(n, len(h.Sum(nil))) // 4 32
```

## Your task

Make `checkDelivery` check in this order, and return the receipt only when every check passes:

```text
decodeEnvelope          ErrBadEnvelope
findKey(keys, e.KeyID)  ErrUnknownSigner
the key not revoked     ErrRevoked
openReceipt             ErrBadSignature, ...
the receipt's KeyID     ErrBadEnvelope
SignedAt in the window  ErrOutsideWindow
checkFile               ErrWrongFile
```

Then make `checkFile` hash the file with the receipt's own hash (`parseDigest`, `newHash`) and return `ErrWrongFile` unless its size and digest are the receipt's.
