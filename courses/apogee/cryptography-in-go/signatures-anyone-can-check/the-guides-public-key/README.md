---
slug: the-guides-public-key
title: The Guide's Public Key
type: code_tests
position: 3
difficulty: 8
topics: crypto, interfaces
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Tuesday afternoon. Yusra's pages for customers' engineers, the customer guide, must hand every customer the public key that checks Apogee's receipts, in a form their own tools can read, before the first receipt is signed on Friday.

Tools everywhere read public keys as PEM: base64 between a `BEGIN PUBLIC KEY` line and an `END` line. `x509.MarshalPKIXPublicKey` writes a key in the standard form that goes inside, and `pem.EncodeToMemory` wraps it. Reading back, `pem.Decode` returns the first block and the text after it, and `x509.ParsePKIXPublicKey` returns an `any` that may hold any kind of key. A type assertion with `ok` asks which:

```go
var v any = "YST"
code, ok := v.(string)
fmt.Println(code, ok) // YST true
_, ok = v.(int)
fmt.Println(ok) // false
```

## Your task

1. Make `publishKey` return `pub` as a `PUBLIC KEY` block.
2. Make `readKey` refuse input whose first block isn't a `PUBLIC KEY` block or is followed by anything but white space, and return the key inside if it is an `ed25519.PublicKey`, else the error `not an Ed25519 key`.
