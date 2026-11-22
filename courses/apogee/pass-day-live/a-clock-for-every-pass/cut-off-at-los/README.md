---
slug: cut-off-at-los
title: Cut Off at LOS
type: code_tests
position: 2
difficulty: 9
topics: context, networking
created: 2026-11-11T16:41:21Z
modified: 2026-11-12T02:35:46Z
---
Thursday 5 July, first thing. On Wednesday's run the rehearsal hung for a simulated week behind Calidra Atoll, whose host's router answers a call and never says a word.

A read on a connection with no deadline returns only when bytes come or the connection is closed; a context that ends wakes nothing by itself. So a pass still waiting for a banner outlives its LOS and holds its dish, and the week never ends.

## Your task

End every call at its pass's LOS, whatever it is waiting for. A pass whose station never sends its banner is hung up at its LOS, and its LOS line ends:

```text
failed no_answer
```

Every other pass flies on, and the week ends at its last LOS. The tests make Calidra Atoll silent for a whole week beside Port Ellery, and check that each Calidra call is hung up at exactly its LOS. A run still going an hour after its week is stopped by the tests and reported as still running.
