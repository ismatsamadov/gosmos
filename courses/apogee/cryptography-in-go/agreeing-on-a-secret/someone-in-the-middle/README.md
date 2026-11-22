---
slug: someone-in-the-middle
title: Someone in the Middle
type: code_tests
position: 3
difficulty: 9
topics: crypto, security
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T07:46:46Z
---
Tuesday afternoon, Mirela plays the company that runs the backhaul. She answers each end with a share of her own and agrees on a key with each, while both ends believe they agreed with each other: she is the machine in the middle.

A share needs proof of who sent it. Each station also keeps a lasting Ed25519 key, its identity key, whose public half the other end holds in advance, and signs its share with it. One line carries both:

```text
signed: backhaul share|PEL|a10c…a93e
line:   HELLO PEL a10c…a93e 6da1…390a
```

The middle can make shares, but no signature from PEL's or YST's key.

## Your task

1. Make `hello` sign those bytes with `id.Key` and append the signature as 128 hex digits.
2. Make `checkHello` refuse a line without exactly four fields, a code not in `known`, a bad share and a signature that fails, each with an error that wraps `ErrBadHello` and reads like these:

```text
bad hello: unknown station "MID"
bad hello: signature does not verify
```

> [!NOTE]
> `known` holds the other ends' keys, never a station's own, so a HELLO sent back to its sender is refused too.
