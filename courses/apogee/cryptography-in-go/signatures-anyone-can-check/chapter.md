---
slug: signatures-anyone-can-check
title: Signatures Anyone Can Check
position: 6
---
Sign and verify with `crypto/ed25519`, publish public keys as PEM through `x509.MarshalPKIXPublicKey`, sign the exact bytes you send with `json.RawMessage`, choose verifying keys by ID, and check a delivered file against its receipt.
