---
slug: sealed-for-cirrine
title: Sealed for Cirrine
type: code_tests
position: 2
difficulty: 9
topics: crypto, security
created: 2026-11-12T02:35:46Z
modified: 2026-11-12T11:16:23Z
---
Tuesday evening, and every Cirrine pass in the rehearsal still fails refused. Cirrine's contract says nobody but Cirrine and Apogee may read its commands, now those of thirteen satellites.

A sealed command hides `<OP>|<args...>` inside AES-256-GCM under the satellite's seal key (info `seal|<sat>`), with a fresh 12-byte nonce and the header, `SEAL|2|<sat>|<seq>|<key ID>`, as associated data: anyone can read the header, nobody can change it. With the nonce `a1b2…4b5c`, Cirrine-1's first PING would read, on one line:

```text
SEAL 2 cirrine-1 9713 k-2028-10
    a1b2…4b5c fc5a…379c
```

The key ID names the master. Then come the nonce and the sealed bytes, in lowercase hex.

## Your task

Send every command to a Cirrine Atmospherics satellite as a SEAL line, under the master its signature would use; every other satellite's PING stays a CMD line.

Draw each nonce from `crypto/rand`: no clock or seed can make it again.

The tests fly the week: every Cirrine command opened and run, Cirrine-1's first line checked against the one above with its own nonce, and every nonce from `crypto/rand`, none sealing two lines in that run or another.
