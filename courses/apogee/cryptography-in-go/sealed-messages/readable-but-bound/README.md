---
slug: readable-but-bound
title: Readable but Bound
type: code_tests
position: 4
difficulty: 8
topics: crypto, security
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Monday 21 August. Ystmark's uplinkd queues each command for its pass, so it must read a sealed command's satellite and number, never the command. Mirela gives an old sealed command a new number, to slip it past the check that refuses any number not above the last, the replay guard.

The last argument of `Seal` and `Open` is associated data: bytes the tag covers but `Seal` leaves out of its result, so they travel readable beside the sealed ones. Open with any other associated data, and `Open` fails:

```go
note := []byte("coffee low")
sealed := aead.Seal(nil, n, note, []byte("YST"))
_, err := aead.Open(nil, n, sealed, []byte("PEL"))
fmt.Println(err != nil) // true: another label
```

`Sealed` holds the readable header, and `s.AD()` writes it as associated data: `SEAL|2|cirrine-1|4471|k-2028-07`.

## Your task

Pass `s.AD()` as the associated data in `sealWithHeader` and in `openWithHeader`, which both pass `nil` now. The test moves your sealed bytes to a new number, satellite, key and version.
