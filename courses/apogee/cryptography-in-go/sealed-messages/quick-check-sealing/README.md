---
slug: quick-check-sealing
title: Quick Check: Sealing
type: quiz
position: 3
difficulty: 7
topics: crypto, security
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Thursday. Before the first sealed pass, Nkem reads the week's sealing code and draws its path on the whiteboard, from uplinkd to the satellite:

```go
sealed := aead.Seal(nil, nonce, cmd, nil)
// on the air: the nonce, then sealed
got, err := aead.Open(nil, nonce, sealed, nil)
```

Then she asks her questions one at a time, as she always does.

## Your task

Answer Nkem's three questions: what `Open` does with a changed byte, where a nonce may travel, and what `Open` returns with the right key and the wrong nonce.
