---
slug: never-the-raw-secret
title: Never the Raw Secret
type: code_tests
position: 4
difficulty: 8
topics: crypto, security
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T07:46:46Z
---
Wednesday. The draft seals backhaul traffic with the X25519 secret itself as its AES key, and Mirela stops it there: "Thirty-two bytes, but not thirty-two random bytes. And nothing in them says which session they belong to."

An X25519 secret is a coordinate of a point on a curve, with structure a key shouldn't have, so it goes through HKDF first, as the satellites' keys did in August. The result is the session key. HKDF's info ties it to this session's two shares, the dialer's first (Port Ellery dials), so ends that saw different shares never share a key:

```text
apogee backhaul v1|<dialer hex>|<answerer hex>
```

```go
a, _ := hkdf.Key(sha256.New, s, nil, "a", 32)
b, _ := hkdf.Key(sha256.New, s, nil, "b", 32)
fmt.Println(bytes.Equal(a, b)) // false
```

## Your task

Make `sessionKey` return an error for a secret that isn't 32 bytes, then return 32 bytes from `hkdf.Key` (package `crypto/hkdf`) with SHA-256, a nil `salt`, and that info.
