---
slug: quick-check-which-key-fits
title: Quick Check: Which Key Fits
type: quiz
position: 6
difficulty: 9
topics: crypto, security
created: 2026-11-13T03:46:40Z
modified: 2026-11-13T12:17:22Z
---
Friday, 19:20. The handover is ready for the morning, and Nkem lists the season's three kinds of key pair on the whiteboard, one to a row, with RSA at 3,072 bits, about as strong as the other two:

```text
RSA      lock, sign         signature 384 bytes
Ed25519  sign               signature 64 bytes
X25519   agree on a secret  share 32 bytes
```

## Your task

Answer Nkem's three questions, one at a time: which key fits a fresh secret for every backhaul session, why Apogee's receipts stay Ed25519, and why Brinewatch's handover uses RSA at all.
