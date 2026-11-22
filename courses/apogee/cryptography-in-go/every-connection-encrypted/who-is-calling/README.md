---
slug: who-is-calling
title: Who Is Calling
type: code_tests
position: 5
difficulty: 9
topics: crypto, http, security
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T19:06:42Z
---
Both ends are proven now, and passplan, the program that predicts the passes, has a client certificate too, so it can read Ystmark's status. Mirela: "Reading is fine. Queuing a command is uplinkd's job alone."

After a mutual handshake, `r.TLS.VerifiedChains` holds the chains Go verified, each running from the client's own certificate, whose common name says who is calling, up to a root. A request that arrived without TLS has a nil `r.TLS`, so look before you index. On passplan's request, the first chain reads:

```go
for _, c := range r.TLS.VerifiedChains[0] {
	fmt.Println(c.Subject.CommonName)
}
// passplan
// Apogee Stations CA 2028
// Apogee Root CA 2028
```

## Your task

1. Make `caller` return the common name of `r.TLS.VerifiedChains[0][0]` and true, or `""` and false when there is no TLS or no verified chain.
2. At the TODO, answer 403 `only uplinkd may queue commands` with `http.Error` unless `caller` says `uplinkd`.

> [!WARNING]
> Never read the caller from `PeerCertificates`: it holds whatever the client sent.
