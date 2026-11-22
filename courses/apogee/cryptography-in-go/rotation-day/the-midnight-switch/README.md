---
slug: the-midnight-switch
title: The Midnight Switch
type: code_tests
position: 0
difficulty: 9
topics: crypto, time
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Monday 9 October, Mirela's rehearsal of Rotation Day. A command sealed with `k-2028-07` at 23:59:58 reaches Cirrine-1's stand-in at 00:00:03, three seconds after `k-2028-10` takes over. The draft refuses it, and the pass is lost.

A safe switch has an overlap, a stretch of time when both versions are accepted. A satellite accepts a version from the moment it holds it, even before its `Activate`, until its `Retire`, when the version retires. A revoked version, one withdrawn for good, is refused at once. Every window is half open:

```go
// from included, until excluded
in := !t.Before(from) && t.Before(until)
```

## Your task

`accepts` holds the draft's hard switch. Make it accept the version named `id` when `ring` holds it, it isn't `Revoked`, and `now` is before its `Retire`, or its `Retire` is zero.

> [!NOTE]
> `k-2028-07` retires at 2028-10-24T00:00:00Z: an overlap of seven days after `k-2028-10` activates.
