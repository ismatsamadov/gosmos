---
slug: telemetry-in
title: Telemetry In
type: code_tests
position: 1
difficulty: 9
topics: files, maps, mutexes
created: 2026-11-11T22:09:38Z
modified: 2026-11-12T10:46:09Z
---
Friday afternoon. The Board shows each satellite's newest frame, and Joss warns that Thistle-2 still sends `batt_v=nan` when its battery sensor glitches.

A TLM frame carries one telemetry line: a satellite's readings from one moment, numbered by the second it was made. A satellite over two stations at once is heard through both, so its frames can arrive out of number order. `telemetry.go` gives `formatLine` and the strict `parseLine`, which refuses `nan`. The first from `thistle-2` at Port Ellery on Thursday (one line on the wire and in the file):

```text
thistle-2 1879776236 2029-07-26T16:03:56Z
    batt_v=7.66 temp_c=8.4 mode=NOMINAL
```

## Your task

Read every TLM line with `parseLine`, dropping what it refuses, and keep each satellite's frame with the highest number.

Save `reports/health.txt` when the week is done, before the `week done` line: a line per satellite, in the order `newNetwork` lists them, its frame as `formatLine` writes it, or else:

```text
brine-6 never heard
```

A stopped run saves no report. The tests compare the file with each satellite's newest good frame as the simulated world sent it, for the week and for a Thursday when Ystmark and Port Ellery relay every frame late.
