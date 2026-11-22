---
slug: a-certificate-for-ystmark
title: A Certificate for Ystmark
type: code_tests
position: 4
difficulty: 9
topics: crypto, structs
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Wednesday 20 September. On Friday the stations CA signs its first certificate, Ystmark's, serial 1001, with the code you write today. Sunniva: "YST first. Make it count."

A station's certificate is no CA: it names one host and proves one thing. Its common name and its one DNS name are the host; its extended key usage is server authentication, so it proves a server and nothing else; its key usage is digital signature; and it lasts 90 days. The CA's certificate is its parent, and the CA's key signs it. A client's usages, for comparison:

```go
tmpl.ExtKeyUsage = []x509.ExtKeyUsage{
	x509.ExtKeyUsageClientAuth,
}
tmpl.KeyUsage = x509.KeyUsageDigitalSignature
```

## Your task

1. Return an error when `ca.Cert` is nil or `ca.Key` isn't `ed25519.PrivateKeySize` bytes: `x509.CreateCertificate` would panic on either.
2. Fill the template for `host`: serial `serial`, organization `Apogee`, from `now` to 90 days later, and a server's usages.
3. Create it with `x509.CreateCertificate`, parent `ca.Cert` and signer `ca.Key`, and return the parsed certificate. The test checks its chain up to the root.
