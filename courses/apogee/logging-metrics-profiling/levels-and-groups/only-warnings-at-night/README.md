---
slug: only-warnings-at-night
title: Only Warnings at Night
type: code
position: 0
difficulty: 6
topics: logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:37Z
---
Last night the ops console, the web tool where engineers read a station's log, gave up before it showed one of Lantern Atoll's passes. Since Friday's deploy, when the new uplinkd replaced the old one there, it writes every record, the DEBUG ones included, to catch the next failure: 2.1 million lines a night.

A handler writes a record only when the record's level is at least the handler's minimum, which `HandlerOptions.Level` sets; left unset, it is INFO. Levels rise from DEBUG through INFO and WARN to ERROR, so a handler at ERROR drops the other three:

```go
o := &slog.HandlerOptions{
	Level:       slog.LevelError,
	ReplaceAttr: dropTime,
}
l := slog.New(slog.NewTextHandler(os.Stdout, o))
l.Warn("carrier lost", "site", "YST")
l.Error("pass failed", "site", "YST")
// level=ERROR msg="pass failed" site=YST
```

## Your task

The night shift's screen should show only what needs a person. Set the options' `Level` in `main` so the program prints its WARN and ERROR records only:

```text
level=WARN msg="carrier lost" site=LAN
level=ERROR msg="pass failed" site=LAN
```

> [!NOTE]
> `main` still makes all five calls; the handler drops the INFO and DEBUG ones before anything is written.
