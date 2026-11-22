---
slug: both-ends-proven
title: Both Ends Proven
type: code_tests
position: 4
difficulty: 9
topics: crypto, http, security
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T07:46:46Z
---
Ystmark now proves who it is, yet it still answers anyone who connects to it. Mirela: "A client that can't prove it is one of our servers gets nothing."

In mutual TLS, the server asks the client for a certificate too and checks it against roots of its own, `ClientCAs`. The client sends its chain from `Certificates` and proves it holds the key. A certificate also says what it is for, so Ystmark's server certificate can't pass as a client's. A client with nothing to show gets this:

```go
_, err := anon.Get(url + "/status")
fmt.Println(err)
// Get "https://…": remote error:
// tls: certificate required
```

## Your task

1. Make `requireClients` set `ClientAuth` to `tls.RequireAndVerifyClientCert` and `ClientCAs` to `roots`.
2. Make `clientConfig` send `cert` in its `Certificates`.
