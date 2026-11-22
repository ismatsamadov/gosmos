---
slug: name-the-program-once
title: Name the Program Once
type: code_tests
position: 3
difficulty: 6
topics: logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
Tuesday night the night shift searches for `"program":"uplinkd"`, and finds half of uplinkd's lines: those whose call happened to name the program and the site. The other half can't be found.

`Logger.With` returns a new logger that puts the attributes you give it on every record it writes, before the record's own. The logger you called it on doesn't change, so one base logger can make a logger for each program and each place:

```go
ys := base.With("site", "YST")
ys.Info("lines queued", "lines", 2)
// {"level":"INFO","msg":"lines queued",
//  "site":"YST","lines":2}
base.Info("lines queued", "lines", 2)
// {"level":"INFO","msg":"lines queued",
//  "lines":2}
```

## Your task

Make `programLogger` return a logger made from `base` that writes two attributes on every record, `program` and then `site`, with the values it is given.

> [!NOTE]
> `site` is where the program runs: a station's code, such as `LAN`, or `HQ` for the rack at headquarters.
