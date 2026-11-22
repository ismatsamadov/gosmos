---
slug: apogees-own-root
title: Apogee's Own Root
type: code_tests
position: 2
difficulty: 9
topics: crypto, time
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T17:01:10Z
---
Monday 18 September. Mirela makes Apogee's root on a laptop that has never touched a network, then locks laptop and key in the safe: every certificate Apogee trusts will lead back to this one.

A root signs itself, so its template is also its parent. A certificate that may sign others, a certificate authority or CA, says so in its basic constraints, `IsCA` and `BasicConstraintsValid`; `MaxPathLen` 1 lets one more CA stand below it, and its key usage lets it sign certificates and lists of withdrawn ones. A one-day template:

```go
tmpl := &x509.Certificate{
	SerialNumber: big.NewInt(7),
	Subject:      pkix.Name{CommonName: "Test"},
	NotBefore:    day,
	NotAfter:     day.AddDate(0, 0, 1),
}
```

## Your task

1. Return an error for a key that isn't `ed25519.PrivateKeySize` bytes: `key.Public()` panics on a nil one.
2. Fill the template: serial 1, `Apogee Root CA 2028`, organization `Apogee`, `now` to 10 years later, a CA with `MaxPathLen` 1 and `KeyUsageCertSign | KeyUsageCRLSign`.
3. Create it as its own parent, and return `x509.ParseCertificate` of the result.

> [!NOTE]
> Ten years is `now.AddDate(10, 0, 0)`: two leap days fall in those ten years, so 3,650 × 24 hours ends two days short.
