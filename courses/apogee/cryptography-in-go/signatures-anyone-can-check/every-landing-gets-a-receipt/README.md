---
slug: every-landing-gets-a-receipt
title: Every Landing Gets a Receipt
type: code_tests
position: 6
difficulty: 8
topics: crypto, json, time
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Thursday 31 August. From tomorrow, the service customers book passes through, the booking service, signs a receipt for every file delivered, and Sunday's Furrow-A file gets the first.

A receipt takes its satellite and station from the booking, and its file, size and landing time from the landing, the note the delivery program, dlinkd, sends when a file is complete. Its signing time comes from the booking service's clock, to the second. `Truncate` drops whatever is below a unit:

```go
d := 1731 * time.Millisecond
fmt.Println(d.Truncate(time.Second)) // 1s
```

## Your task

Write `issue`. Return an error, and sign nothing, for a landing of another booking or a digest `parseDigest` refuses. Otherwise fill a version 1 `Receipt` with the notary's `KeyID`, the booking's and the landing's fields, the digest and `SignedAt` from `now()` truncated to the second, and return `wrapReceipt`'s envelope.
