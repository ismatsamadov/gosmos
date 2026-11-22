---
slug: lines-for-machines
title: Lines for Machines
type: code_tests
position: 1
difficulty: 6
topics: json, logging
created: 2026-11-07T07:51:36Z
modified: 2026-11-07T07:51:36Z
---
Wednesday, after the 08:00 meeting, Nkem: "`error` is not a log line." The people who run the passes overnight, the night shift, search logs with a tool that reads one JSON object per line, and Lantern Atoll's new lines must go straight into it.

A logger's handler decides how its records look. The text handler writes `key=value` pairs for people; `slog.NewJSONHandler` writes each record as one JSON object on one line, so a program reads the level, the message and every attribute back by name. One record, both ways:

```text
level=INFO msg="command sealed" seq=4688
{"level":"INFO","msg":"command sealed","seq":4688}
```

## Your task

`newLogger` builds its logger on the text handler. Make it use `slog.NewJSONHandler` on the same writer, with the same options, so `dropTime` still leaves the time out.

> [!NOTE]
> A line break inside a value can't split a JSON line: the handler writes it as `\n`, so one record is always one line.
