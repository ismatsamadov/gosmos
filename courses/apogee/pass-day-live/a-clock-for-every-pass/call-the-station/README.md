---
slug: call-the-station
title: Call the Station
type: code_tests
position: 1
difficulty: 8
topics: networking, streams, strings
created: 2026-11-11T16:41:21Z
modified: 2026-11-12T02:35:45Z
---
Wednesday afternoon. Sunniva wants the rehearsal's calls to reach each station's radio on time and leave on time, as the real ones will.

A pass's call to its station is its pass line, and `w.Dial(ctx, code)` makes it. The radio answers with one line, its banner, ended by CR LF or by LF alone, and sometimes in pieces:

```text
PEL READY PROTO 1
```

## Your task

At each contact's AOS, call its station and read the banner: it must name the pass's station and protocol 1. Hold the call until LOS, then hang up. A failed call, or any other banner, fails the pass, and its LOS line, still at LOS, ends:

```text
failed no_answer
```

The simulated world keeps a record of every call: the station, when you dialed and when you hung up. The tests check one call per contact, dialed at its AOS and, once answered, hung up at its LOS. They send Port Ellery's banner a byte at a time, give Ilha Brumosa's radio a banner of protocol 2, and fail every call to Corith Point.
