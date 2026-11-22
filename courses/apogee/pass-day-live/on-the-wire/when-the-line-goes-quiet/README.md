---
slug: when-the-line-goes-quiet
title: When the Line Goes Quiet
type: code_tests
position: 3
difficulty: 10
topics: networking, streams, time
created: 2026-11-11T22:09:38Z
modified: 2026-11-12T02:35:46Z
---
Monday afternoon. Sunniva, on the link: "Storms kill the fiber without a sound. Calls stay open, nothing comes. Backhaul up in ten minutes."

A dead line raises no error: a read on it waits. A read deadline ends the wait, and `conn.SetReadDeadline` takes an instant on the `time` package's own clock, never `w.Now()`'s. Moved forward before every read, it lets a slow frame trickle in and still catches a silent line.

## Your task

Give every read 90 s to bring a byte, or end the call. Until LOS, call a lost line again at once, with HELLO for the same pass and dish. After a call that fails or ends without a WELCOME, wait before the next: 1 s, doubled for each such failure in a row, at most 30 s, times a random factor from 0.5 up to 1. A stop ends the wait at once.

The tests kill Ystmark's fiber at 03:31 on Thursday, during a Pipit-2 pass: calls fail until the backhaul is up at 03:41, then frames trickle in, 64 bytes at a time. Calidra Atoll's and Cape Wren's calls, which never get a WELCOME, now end after 90 s, and the next waits.
