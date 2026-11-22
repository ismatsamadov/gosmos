---
slug: signed-twice-never-alike
title: Signed Twice, Never Alike
type: code_tests
position: 5
difficulty: 9
topics: crypto, security
created: 2026-11-13T03:46:39Z
modified: 2026-11-13T03:46:39Z
---
Brinewatch's security office will approve tomorrow's handover with one line signed by its new signing key, and Yusra has just read that key's whole fingerprint, all 64 digits, back to Brinewatch by phone. Apogee locks nothing until the signature checks out, and to try it tonight Mirela plays Brinewatch, so you write both sides.

RSA signs a digest, not the message: hash the request with SHA-256, then sign its 32 bytes and name their hash, `crypto.SHA256`. Its safe form, PSS, mixes fresh random bytes into every signature, so unlike Ed25519 it never signs the same request the same way twice, and a check accepts each one:

```go
// k: a private key, h: crypto.SHA256, d: a digest
a, _ := rsa.SignPSS(rand.Reader, k, h, d, nil)
b, _ := rsa.SignPSS(rand.Reader, k, h, d, nil)
fmt.Println(bytes.Equal(a, b)) // false
```

## Your task

1. Make `signRequest` sign `req` with `rsa.SignPSS`, SHA-256 and `nil` options.
2. Make `checkRequest` return what `rsa.VerifyPSS` says about `sig` over `req`'s digest, with `nil` options too.
