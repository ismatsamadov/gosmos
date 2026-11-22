---
slug: twelve-bytes-any-bytes
title: Twelve Bytes, Any Bytes
type: code_tests
position: 5
difficulty: 9
topics: bytes, errors, testing
created: 2026-11-11T12:51:44Z
modified: 2026-11-12T08:43:35Z
---
Pipit beacons, the twelve-byte status messages each Pipit satellite sends, now reach tlmd from the new stations' receivers too, and when a pass ends in the middle of one, a receiver hands over only its first few bytes.

A fuzz target can take a `[]byte` as well as a string, and a round trip is a property that needs no answer: decode the bytes, encode the result, and the same bytes must come back. A corpus file holds bytes as it holds a line:

```text
go test fuzz v1
[]byte("AP\x01")
```

## Your task

`decodeBeacon`, in `decode.go`, reads `b` before it knows its length, so the 7 bytes saved in `testdata/fuzz/FuzzBeacon/` make it panic. Before anything else, refuse any length but 12, with the length in this error:

```text
beacon: want 12 bytes, got 7
```

> [!TIP]
> Thirteen bytes are wrong too: `FuzzBeacon` encodes every beacon it accepts and wants back exactly the bytes it gave.
