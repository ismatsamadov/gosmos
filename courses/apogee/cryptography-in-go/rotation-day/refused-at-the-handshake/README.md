---
slug: refused-at-the-handshake
title: Refused at the Handshake
type: code_tests
position: 5
difficulty: 10
topics: crypto, networking, security
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
10:14. List 1 is signed, and Mirela wants uplinkd to refuse certificate 1005 from its very next connection: "Its chain still verifies. That's exactly why we check the list."

`tls.Config.VerifyConnection` runs during every handshake, a resumed session's too, after the usual certificate checks, with the `ConnectionState` so far; an error from it ends the handshake. On the client's side, the station's own certificate is `PeerCertificates[0]`:

```go
check := func(cs tls.ConnectionState) error {
	log.Println(cs.ServerName)
	return nil // nil lets the handshake go on
}
cfg.VerifyConnection = check
```

## Your task

Make `withRevocations` set `cfg.VerifyConnection` to check that certificate with `isRevoked` against `current()` at `now()`, both called during the handshake. Return its error, or the refusal below, with the certificate's own serial; until the first list is out, `current()` returns nil: refuse nothing then.

```text
certificate 1005 is revoked
```
