---
slug: never-the-key-it-brings
title: Never the Key It Brings
type: code_tests
position: 5
difficulty: 8
topics: crypto, json, maps
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Wednesday afternoon. A draft envelope carries its own public key, so customers would need nothing from the guide. Mirela signs a receipt with a key of her own, puts that key in the envelope, and the draft accepts it.

A key that comes with the message proves only that someone made a key. The key that checks a receipt comes from your own list, picked by the envelope's `key_id`, and an ID that isn't there is refused. A map lookup with `ok` tells a missing key from a present one:

```go
names := map[string]string{"YST": "Ystmark"}
name, ok := names["XYZ"]
fmt.Printf("%q %t\n", name, ok) // "" false
```

## Your task

1. Decode the envelope with `decodeEnvelope`, which refuses any field `Envelope` doesn't have, such as `pub`.
2. Take the key from `keys` by `e.KeyID`, wrapping `ErrUnknownSigner` when there is none, and open the receipt with `openReceipt`.
3. Wrap `ErrBadEnvelope` when the receipt's own `KeyID` differs from the envelope's.
