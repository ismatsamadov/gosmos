---
slug: ystmarks-own-certificate
title: Ystmark's Own Certificate
type: code_tests
position: 1
difficulty: 9
topics: crypto, http, time
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T07:46:46Z
---
Ystmark's uplinkd starts serving with the certificate Apogee's certificate authority issued it on Friday 22 September. Every client holds Apogee's root and nothing else, so the station must send the middle of its chain as well.

A server lists what it presents in its config's `Certificates`. Each `tls.Certificate` holds a chain as DER bytes, every `x509.Certificate`'s `Raw`: the server's own certificate first, then the intermediate that signed it, with the private key of the first. A client reads the chain it was sent:

```go
for _, c := range resp.TLS.PeerCertificates {
	name := c.Subject.CommonName
	fmt.Println(c.SerialNumber, name)
}
// 1001 yst.apogee.internal
// 2 Apogee Stations CA 2028
```

## Your task

Make `serverConfig` present one `tls.Certificate`: `leaf.Raw`, then `inter.Raw`, with `key` as its `PrivateKey` and `leaf` as its `Leaf`. Set `Time: now` too, so the station judges dates by the test's clock.
