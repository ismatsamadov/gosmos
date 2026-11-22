---
slug: ystmarks-storm-key
title: Ystmark's Storm Key
type: code_tests
position: 7
difficulty: 9
topics: crypto, networking
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T07:46:46Z
---
Thursday 7 September. The fiber is still cut, and tonight this week's command file for Ystmark crosses the backhaul, sealed under a key the two ends agree on as the session starts. Sunniva: "Wind 40, dropping. Ready when you are."

The handshake puts the week's pieces in order. Port Ellery dials and sends its HELLO first, with a fresh share; Ystmark reads it with `readHello`, which checks it, and only then answers with a HELLO of its own. Each end derives the session key with `sessionKey`, the dialer's share first:

```text
PEL -> YST  HELLO PEL a10c…a93e 6da1…390a
YST -> PEL  HELLO YST b6f8…5a66 ebe7…6203
```

## Your task

Write `dial`, Port Ellery's side, and `answer`, Ystmark's: each makes a fresh X25519 key, sends its HELLO with `writeLine`, checks the other's with `readHello`, and returns the session key. `dial` also refuses an answer from any code but `peer`, with an error wrapping `ErrBadHello`; `answer` returns the dialer's code. The given `sendData` and `readData` then carry the file.

> [!TIP]
> Over `net.Pipe`, a write waits for its reader: if both ends write first, neither ever reads.
