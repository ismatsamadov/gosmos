---
slug: the-night-shift-runs-a-pass
title: The Night Shift Runs a Pass
type: code_tests
position: 3
difficulty: 8
topics: files, testing
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T16:59:09Z
---
Friday night, 04:00 in the Fishbowl, the glass-walled ops room. Ingrid runs the night's passes from your prompt, and the forecast says Cape Wren's wind will reach 72 km/h by 04:40.

The tests now replay whole sessions from `testdata`. Each `.in` file is what the night shift types, with lines such as `@04:40` moving the clock between commands, and each `.golden` file is the screen it must produce, word for word.

## Your task

Two things are missing for tonight:

1. Make `next` skip passes at stowed stations.
2. End a `passes` line with `(stowed)`, after a space, when its station is stowed.

Then every session in `testdata` matches its transcript, including:

```text
gctl> passes CWR
04:58-05:06 CWR furrow-a max 33.5 (stowed)
gctl> next
next: brine-3 over LAN at 05:20, in 40m0s
```
