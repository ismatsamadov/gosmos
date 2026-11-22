---
slug: hello-in-frames
title: Hello in Frames
type: code_tests
position: 0
difficulty: 9
topics: bytes, networking, streams
created: 2026-11-11T22:09:38Z
modified: 2026-11-12T10:46:08Z
---
Sunniva's news: the radios now want each call to name its pass and dish first.

After the banner come frames: length, type, seq, body and CRC, every number big-endian (`passline.go` names the types). HELLO for Thistle-2 at Port Ellery, Thursday:

```text
00 1e            length: 30 bytes follow
01               type: HELLO
00 00 00 00      seq: 0
74 68 ... 20 31  body: thistle-2/PEL/02136 1
ba 2c bf 54      CRC-32 of type, seq, body
```

## Your task

Send HELLO with the pass and its dish, then read every frame until LOS. The station's refusal, a NAK, carries a code. The LOS line ends with the first of these that fits, or else `flown`: a NAK with code 2, the dishes stowed in the wind; any other NAK; no WELCOME; a BYE `stowed`, cutting the pass short:

```text
skipped stowed
failed refused
failed no_answer
failed stowed
```

Hang up at LOS, or when the station refuses the pass or says BYE. The tests check HELLO's bytes, a banner and WELCOME in one write, Ystmark's wind lock, a radio that never welcomes, and an archive missing a pass, so the next gets the wrong dish.
