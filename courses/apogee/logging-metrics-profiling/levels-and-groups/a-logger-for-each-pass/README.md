---
slug: a-logger-for-each-pass
title: A Logger for Each Pass
type: code_tests
position: 5
difficulty: 7
topics: logging
created: 2026-11-10T15:36:37Z
modified: 2026-11-11T17:33:37Z
---
Wednesday afternoon: uplinkd's `runPass` writes six records for every pass, from `pass started` to `pass flown`, and each has to say which pass it belongs to; on 8 November no line did. With its dish and its smaller antennas, the station can fly two passes at once, so their lines interleave.

`With` takes `slog.Attr` values too, a group among them. A logger made for one piece of work carries its name on every record, and the code doing the work never has to repeat it. Each piece of work gets its own logger, so two running at once never mix:

```go
id := slog.String("booking_id", "bk-5611")
bl := base.With(id)
bl.Info("booking created")
bl.Info("pass planned")
// both lines: ...,"booking_id":"bk-5611"}
```

## Your task

Make `passLogger` return a logger made from `l` that writes two attributes on every record: `pass_id`, the pass's ID, then the group from your `passGroup`, beside it in `levels.go`.

> [!NOTE]
> `runPass` gets its logger once, from `passLogger`, and hands it to every step of the pass.
