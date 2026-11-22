---
slug: certificates-for-every-station
title: Certificates for Every Station
type: code_tests
position: 7
difficulty: 9
topics: crypto, files
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T19:06:42Z
---
Friday 22 September. Before the weekend's passes, every station, uplinkd, passplan, the program that plans the passes, and dishd, which drives Ystmark's dishes, get a certificate, each for a new key of its own, and Mirela wants every key file readable by its owner alone.

Each holder gets two files. `<name>.pem` holds its certificate, then the stations CA's, so whoever checks it has the whole chain; `<name>.key` holds its private key as a PKCS #8 `PRIVATE KEY` block, with mode 0600. A PEM block is a type and the bytes it wraps:

```go
b := &pem.Block{Type: "PUBLIC KEY", Bytes: der}
os.Stdout.Write(pem.EncodeToMemory(b))
```

## Your task

1. In `issueFiles`, make a key with `ed25519.GenerateKey`, issue the holder's certificate with `issueServer`, or `issueClient` for a client, and encode both files, the key through `x509.MarshalPKCS8PrivateKey`.
2. In `issueAll`, do that for every entry of `holders` first, then save each file with `saveFile`, keys with mode 0600, so a failure leaves no files at all.
