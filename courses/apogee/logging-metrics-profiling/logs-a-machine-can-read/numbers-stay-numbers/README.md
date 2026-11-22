---
slug: numbers-stay-numbers
title: Numbers Stay Numbers
type: code_tests
position: 3
difficulty: 6
topics: json, logging
created: 2026-11-07T07:51:36Z
modified: 2026-11-07T07:51:36Z
---
Tomás suspects Lantern Atoll's diesel generator and wants every failed pass that climbed above 30° of elevation, the satellite's angle above the horizon: "Downpours and palm trees only take the low ones!" His search compares numbers, and a number written as text can't be compared with 30.

`slog.String`, `slog.Int` and `slog.Float64` each make one attribute of their type, and `Logger.LogAttrs` writes a record from a context, a level, a message and any number of attributes. A float stays a JSON number:

```go
ctx := context.Background()
l.LogAttrs(ctx, slog.LevelWarn, "carrier lost",
	slog.String("site", "LAN"),
	slog.Float64("margin_db", -1.5))
// {"level":"WARN","msg":"carrier lost",
//  "site":"LAN","margin_db":-1.5}
```

## Your task

Make `logFlown` write one INFO record with `l.LogAttrs`: the message `pass flown` and three typed attributes, `pass_id`, `max_elev_deg` and `lines`.

> [!NOTE]
> JSON writes a whole number of degrees without its point: 8.0 comes out as `8`, still a number.
