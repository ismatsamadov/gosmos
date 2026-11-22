---
slug: copies-not-references
title: Copies, Not References
type: code
position: 3
difficulty: 4
topics: structs, variables
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Ingrid wants to see what a 12° mask would do to Cape Wren's passes, but the live record feeds the scheduler. "Try it on a draft. Don't touch live."

Assigning a struct copies every field. From then on the two variables are separate values, so setting a field of one leaves the other alone:

```go
plan := Orbit{AltKm: 525, IncDeg: 97.5}
backup := plan
plan.AltKm = 510
fmt.Println(backup.AltKm) // 525
```

## Your task

At the `TODO`, set the draft's mask to 12 degrees. The live record must keep its 5.
