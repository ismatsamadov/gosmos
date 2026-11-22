---
slug: fresh-keys-every-session
title: Fresh Keys Every Session
type: code_tests
position: 5
difficulty: 8
topics: crypto, methods
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T07:46:46Z
---
Mirela reads the draft's config: Ystmark makes one X25519 key at its first session and keeps it for every session after. "Record the backhaul today, steal that file next year, and you open every session it ever served."

A new X25519 key for every session, kept nowhere once the session ends, gives forward secrecy: whoever takes the station later finds only its identity key, which could sign shares for new sessions until it is replaced, but opens no recording. One key kept for good turns every recording into a target:

```go
// a year later, with the stolen key:
secret, _ := stolen.ECDH(recordedShare)
// the secret of every session it served
```

## Your task

Make `NewSession` return a new key from `ecdh.X25519().GenerateKey` on every call, and remove the `static` field, so a `Backhaul` keeps nothing. The test copies a `Backhaul` after a session, as a thief would, and tries to open the session it recorded.
