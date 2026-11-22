---
slug: every-command-signed
title: Every Command Signed
type: code_tests
position: 0
difficulty: 9
topics: crypto, security, select
created: 2026-11-12T02:35:46Z
modified: 2026-11-12T11:16:23Z
---
Mirela, before the rehearsal sends its first command: "How will a satellite know it came from us?"

A satellite runs a command only when its sequence number is above the last one it ran, which the WELCOME carries as `last`, and its signature checks out under the satellite's sign key: a derived key, 32 bytes of HKDF-SHA256 from a master key in `uplink-keys.txt`, with salt `apogee uplink keys v1` and info `sign|<sat>`. Pipit-2's first PING, which only asks for an answer, ends in the HMAC-SHA256 of `pipit-2|6315|PING` under its key from `k-2028-10`:

```text
CMD pipit-2 6315 PING sig=0cb1…4810
```

## Your task

Read `uplink-keys.txt`. On every contact, 10 s after AOS once welcomed, send a PING in an UPLINK frame, signed with the newest master its satellite holds: at the window's start it holds every one activated before then that isn't retired or revoked. Keep reading frames while you wait for its ACK. A NAK comes right after `no_answer` in the outcomes, and no answer by LOS last:

```text
failed refused
failed no_ack
```

Cirrine's satellites refuse every CMD line for now. The tests fly Tuesday at Port Ellery, Thistle-2 from `last` 0, and check Pipit-2's first line.
