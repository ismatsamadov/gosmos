---
slug: signed-by-the-ground
title: Signed by the Ground
position: 2
---
Authenticate with `crypto/hmac` rather than a hash of key and message, make canonical forms unambiguous, verify before changing state, compare with `hmac.Equal` and `subtle.ConstantTimeCompare`, and give every message its own HMAC.
