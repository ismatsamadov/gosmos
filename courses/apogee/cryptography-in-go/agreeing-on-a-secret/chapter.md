---
slug: agreeing-on-a-secret
title: Agreeing on a Secret
position: 7
---
Agree on a shared secret with `ecdh.X25519`, check public shares, sign them with Ed25519 against known keys, derive session keys with HKDF, make fresh keys per session, and encapsulate a key with ML-KEM-768.
