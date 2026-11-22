---
slug: quick-check-mutual-proof
title: Quick Check: Mutual Proof
type: quiz
position: 6
difficulty: 9
topics: http, security
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T07:46:46Z
---
Monday 2 October. Mirela finds a draft of Ystmark's handler from before both ends were proven. Its server asked for client certificates with `tls.RequestClientCert`, which requests one and verifies nothing, and the handler read its caller like this:

```go
peer := r.TLS.PeerCertificates[0]
name := peer.Subject.CommonName
```

## Your task

Answer Mirela's two questions: what that draft would believe about a certificate a client made itself, and who decides what a caller with a verified certificate may do.
