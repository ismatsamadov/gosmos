---
slug: the-middle-of-the-chain
title: The Middle of the Chain
type: code_tests
position: 5
difficulty: 9
topics: crypto, security
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Wednesday afternoon. The root stays in the safe, so a station shows its own certificate together with the stations CA's, and Mirela tries the obvious trick at once: Ystmark's certificate signing one for Port Ellery.

`Verify` builds a path from a certificate up to one of your roots, through the certificates in `Intermediates`. Those are stepping stones, never trusted on their own: each must be a CA, stand where its `MaxPathLen` allows, and be signed by the next one up. A station's certificate is no CA, so nothing it signs leads anywhere. A pool can also be filled from PEM text:

```go
pool := x509.NewCertPool()
ok := pool.AppendCertsFromPEM(chainPEM)
fmt.Println(ok) // true
```

## Your task

Make `verifyStation` add every certificate of `chain` to an `Intermediates` pool, and ask for `ExtKeyUsageServerAuth` alone. The test sends Ystmark's certificate with and without its chain, one that Ystmark's key signed, one under an intermediate from Mirela's own root, and uplinkd's client certificate.
