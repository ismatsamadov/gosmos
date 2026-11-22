---
slug: every-contact-of-the-week
title: Every Contact of the Week
type: code_tests
position: 1
difficulty: 8
topics: files, iterators, slices, sorting
created: 2026-11-10T20:36:24Z
modified: 2026-11-12T02:35:45Z
---
Monday afternoon. Each station's dish loader, the program that hands its dishes their passes, wants the week's plan a day per file, and Nkem wants it to be passplan 2's own, contact for contact.

`passes.txt` is passplan's pass archive: a line per pass at the rehearsed stations, in order of AOS, the moment the satellite rises. Each holds the pass ID, AOS, LOS and maximum elevation (one line in the file):

```text
thistle-2/PEL/02136 2029-07-26T16:03:46Z
    2029-07-26T16:09:59Z 11.7
```

## Your task

1. Keep the candidates: passes whose station can serve the satellite (`canServe`, in `plan.go`) and that climb at least 10°. A bad line is an error, `archive line 3: ...`, as is a missing file, and nothing is printed.
2. Give each station's candidates dishes with `assign`, every dish free at the start.
3. For each day, save `plan/<day>.txt` through a temporary file and a rename: the contacts whose AOS falls that day, by AOS then pass ID, a line each with AOS, LOS, dish and pass ID:

```text
16:03:46 16:09:59 1 thistle-2/PEL/02136
```

Then print the totals:

```text
plan: 14297 contacts, 28680 dropped
```

> [!NOTE]
> The tests compare every file with passplan 2's plan, line for line.
