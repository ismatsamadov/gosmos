---
slug: a-pad-used-twice
title: A Pad Used Twice
type: code_tests
position: 0
difficulty: 9
topics: bytes, crypto, security
created: 2026-11-13T03:46:39Z
modified: 2026-11-13T03:46:39Z
---
Friday 6 October. Brinewatch wants Rotation Day's new secrets the way it got its first two, in November and January: XORed with 256 random bytes on a USB stick its courier brought, one copy at each end, its pad. Mirela lays the two old emails side by side. Both used the pad from its first byte.

XOR (`^`) flips a message's bits wherever the pad has a 1, and the same pad bytes flip them back. Random pad bytes, as many as the message's and never used again, give away nothing but its length: a one-time pad. Used twice, it cancels out, as August's repeated nonce did: the two sealed messages XORed together are the two messages XORed together.

```go
m, p := byte('K'), byte(0xa7)
c := m ^ p
fmt.Printf("%#x %c\n", c, c^p) // 0xec K
```

## Your task

`Seal`, Brinewatch's tool, XORs every message with the pad from its first byte. Make it start at the first unused byte, add what it used to `p.used`, and refuse a message longer than what's left, spending nothing:

```text
pad has 7 bytes left, need 25
```
