---
slug: never-twice
title: Never Twice
type: code_tests
position: 3
difficulty: 10
topics: errors, networking, security
created: 2026-11-12T02:35:46Z
modified: 2026-11-12T11:16:23Z
---
Wednesday 11 July. In last night's run, Ystmark's line dropped between a PING to Cirrine-7 and its answer, and the pass failed `no_ack`, though the satellite had run it. Sent again under a new number, it would have run twice.

When the line comes back, the new WELCOME's `last` settles the command in flight. At or above its number, the command ran: count it answered. Below, it hasn't run: send the same line again, same number, same bytes; a sealed one isn't sealed anew. A satellite can still run a command that arrived before the line dropped, so a NAK `replayed sequence number` to the copy means the first ran.

## Your task

After a lost line, settle any command in flight from the next WELCOME's `last`, sending it again only when it hasn't run, and count a copy refused as replayed as run.

The tests drop last night's call at Ystmark right after its PING runs, then right before, and lose Cirrine-1's line over the backhaul while its PING is on the way: each pass flown, each PING run once, and a dropped one sent again as it was.
