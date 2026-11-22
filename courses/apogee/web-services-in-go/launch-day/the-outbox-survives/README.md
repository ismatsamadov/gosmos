---
slug: the-outbox-survives
title: The Outbox Survives
type: code_tests
position: 6
difficulty: 9
topics: errors, files, json
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Thursday's rehearsal restarted the service during Furrow-A's pass, as Ingrid's checklist says, while an event for Cirrine was still being retried. Its retries died with the old process, and nobody would have known.

The outbox keeps every delivery not yet made, and the number of the next event, in `outbox.json`, saved whole after every change through `replaceFile`. A new process opens it and goes on where the old one stopped, with each delivery's attempts and an event count that never repeats:

```json
{"next_event":118,"pending":[
 {"customer":"cirrine","event":{…},
  "attempts":2,"next_at":"…"}]}
```

## Your task

1. Write `save`: marshal an `outboxFile` holding `o.next` and `o.pending`, and write it with `replaceFile`.
2. Make `OpenOutbox` load the file into the outbox it returns. A missing file is the empty outbox it returns now; a file that doesn't decode is an error.
