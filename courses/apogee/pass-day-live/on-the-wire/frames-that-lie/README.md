---
slug: frames-that-lie
title: Frames That Lie
type: code_tests
position: 2
difficulty: 9
topics: bytes, errors, networking
created: 2026-11-11T22:09:38Z
modified: 2026-11-12T10:46:09Z
---
Monday 9 July. Varsk Island's host line flips a bit now and then, and on Friday one damaged frame left a pass deaf until its LOS.

A frame whose CRC doesn't match was damaged on the way, but its length still says where it ends: drop it, and the next frame starts in the right place. A length under 9 or over 1024 means the stream is out of step, and nothing after it can be trusted: the line is lost.

## Your task

Drop every frame whose CRC is wrong, and read on. When a length is out of range, or the call ends before LOS without a BYE, the line is lost: hang up, and unless one of the four earlier outcomes fits, the LOS line ends:

```text
failed line_lost
```

Dialing again comes next lesson. The tests flip a byte in the third TLM frame of a Varsk Island pass, and check that every frame after it is read and the pass flown. They also send a length of 8 right after another pass's WELCOME.
