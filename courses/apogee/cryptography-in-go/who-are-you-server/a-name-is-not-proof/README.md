---
slug: a-name-is-not-proof
title: A Name Is Not Proof
type: code_tests
position: 0
difficulty: 8
topics: crypto, security
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Monday 11 September. Mirela's audit starts with the draft checker for the signed files that tie a station's public key to its name, called certificates: it trusts any certificate whose name says `yst.apogee.internal`. She makes one in a millisecond, signed by a key of her own whose certificate copies the name `Apogee Root CA 2028`.

Anyone can write any name into a certificate, but only a key's holder can sign with that key. A certificate counts only when `Verify` follows its signatures up to a root, a certificate you trust from the start: here, Apogee's own. `Verify` also checks the name, the use the certificate allows, and its dates at `CurrentTime`:

```go
_, err := client.Verify(x509.VerifyOptions{
	Roots: roots,
	KeyUsages: []x509.ExtKeyUsage{
		x509.ExtKeyUsageClientAuth,
	},
})
```

## Your task

Make `trusted` return the error of `cert.Verify` with `roots` as `Roots`, `inters` as `Intermediates`, `host` as `DNSName`, `ExtKeyUsageServerAuth` as the only usage, and `now` as `CurrentTime`.

> [!WARNING]
> Never leave `Roots` nil: `Verify` would ask the Mac's own list of roots, which has never heard of Apogee's.
