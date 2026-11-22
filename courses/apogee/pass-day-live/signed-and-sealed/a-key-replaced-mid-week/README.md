---
slug: a-key-replaced-mid-week
title: A Key Replaced Mid-Week
type: code_tests
position: 4
difficulty: 10
topics: crypto, security, time
created: 2026-11-12T02:35:46Z
modified: 2026-11-12T11:16:23Z
---
Wednesday afternoon. On Monday Mirela made `k-2029-07` to take over at 00:00 on Wednesday 25 July, in the middle of the investors' week, and their reviewer will watch it happen.

A satellite checks commands made with a master only once it holds that master's keys. A LOAD_KEY brings them up, sealed under the newest master it already holds. Thistle-2's, opened, with its sign and seal keys shortened from 64 hex digits each:

```text
LOAD_KEY|k-2029-07|f150…84e0|2b8e…c459
```

## Your task

On each satellite's first contact, after its PING, send a LOAD_KEY for every master in `uplink-keys.txt` it doesn't hold and that is neither retired nor revoked. The satellite holds that master once its LOAD_KEY runs; one unanswered at LOS is settled on its next contact, by the WELCOME's `last`. Sign and seal with the newest master that has activated and the satellite holds.

The tests fly the week: all eighty loaded by Monday 09:27:31, `k-2028-10` until Tuesday 23:59:59 and `k-2029-07` from Wednesday 00:00:00, the retired `k-2028-07` never. Then Wednesday at Port Ellery, where no satellite starts with the new master, and a LOAD_KEY whose answer a dead line loses.
