---
slug: what-a-certificate-says
title: What a Certificate Says
type: code
position: 1
difficulty: 8
topics: crypto, files
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T17:01:09Z
---
Tuesday 12 September, a voice memo from Paz at Salar Quieto, the station on the high plateau: an old test script on a spare PC has been asking the computer that turns Chiquita, the smaller dish, for its position every minute, and nobody knew. "Before anything up here gets a certificate, tell me what one says."

`yst.pem` is Mirela's rehearsal copy of the certificate Ystmark gets on 22 September. It holds two PEM blocks: Ystmark's certificate, then the one that issued it. `pem.Decode` returns the first block and the text after it, and `x509.ParseCertificate` reads a block's bytes:

```go
block, _ := pem.Decode(guidePEM)
fmt.Println(block.Type) // PUBLIC KEY
```

## Your task

Make `main` decode the first block of `yst.pem`, parse it, and print six of its fields: its common name, its issuer's, its serial number, its dates (`time.DateOnly`), its DNS names, and whether it may sign other certificates, which makes it a certificate authority, or CA:

```text
subject: yst.apogee.internal
issuer: Apogee Stations CA 2028
serial: 1001
valid: 2028-09-22 to 2028-12-21
names: [yst.apogee.internal]
CA: false
```
