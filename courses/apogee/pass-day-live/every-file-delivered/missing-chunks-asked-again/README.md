---
slug: missing-chunks-asked-again
title: Missing Chunks, Asked Again
type: code_tests
position: 3
difficulty: 10
topics: maps, networking, slices
created: 2026-11-12T05:36:27Z
modified: 2026-11-12T05:36:27Z
---
Friday afternoon. In the rehearsal, Furrow-K's 07:40 image over Vorrel Island on Tuesday 24 July still comes down three chunks short, and Hollis has read the report: "45 OF 48 IS NOT AN IMAGE."

A satellite keeps what it sends, so its next contact can ask again with a DUMP_FILE command, and the station sends FILE again, then the chunks asked for. The command names the file, then the missing chunk numbers ascending, or the name alone for the whole file after a CRC mismatch:

```text
DUMP_FILE furrow-k-20290724-0740.tif 21 22 23
```

## Your task

On each contact, after the PING and any LOAD_KEY, send a DUMP_FILE for each file of the satellite still incomplete that no pass is bringing down. A pass brings a file down from its FILE, or from its DUMP_FILE's ACK, until the file's END or the pass's LOS. Finish the file from the chunks that come, and deliver it; its event keeps its booking's station.

The tests check Furrow-K's image, asked for on its 08:12:10 contact over Varsk Island and landed at 08:12:30, and Thistle-2's first science file after a wrong CRC, a line dead mid-file and a DUMP_FILE whose answer was lost.
