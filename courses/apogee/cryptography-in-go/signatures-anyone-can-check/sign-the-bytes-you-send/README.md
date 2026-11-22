---
slug: sign-the-bytes-you-send
title: Sign the Bytes You Send
type: code_tests
position: 4
difficulty: 8
topics: crypto, json
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Wednesday 30 August. The first sample checker in Yusra's guide for customers decodes a receipt and encodes it again before checking its signature. Mirela adds `"refund":true` inside one of Apogee's receipts, and the sample accepts it.

A signature covers bytes, not meaning. Decoding drops any field the struct lacks, and encoding again writes Apogee's own bytes back, so the check passes on bytes nobody sent. The envelope that carries a receipt and its signature keeps the receipt as a `json.RawMessage`, the bytes exactly as they came:

```go
var e struct{ N json.RawMessage }
json.Unmarshal([]byte(`{"N": { "a" : 1 }}`), &e)
fmt.Printf("%s\n", e.N) // { "a" : 1 }
```

## Your task

In `openEnvelope`, check the signature over `e.Receipt` with `checkSig` before decoding anything from it. Then decode it with `decodeStrict`, which refuses any field `Receipt` doesn't have, and wrap its error with `ErrBadEnvelope`.

> [!NOTE]
> `checkSig` refuses a `sig` that isn't 128 lowercase hex digits with `ErrBadSignature`, as it refuses a wrong one.
