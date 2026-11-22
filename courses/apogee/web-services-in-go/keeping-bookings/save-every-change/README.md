---
slug: save-every-change
title: Save Every Change
type: code_tests
position: 4
difficulty: 8
topics: files, json
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Yusra's TL;DR: "I retire the spreadsheet when a booking survives the service stopping. Not before."

`FileStore` keeps a `MemStore` and, after every change, saves all of its bookings to `bookings.json` through `replaceFile`, which writes a temporary file beside it and renames it into place, so the file is always whole. When a save fails, `Create` and `Delete` undo their change, so memory and file never disagree:

```go
old := cfg.Station
cfg.Station = "CWR"
if err := cfg.save(); err != nil {
	cfg.Station = old // the file never changed
	return err
}
```

## Your task

`save` does nothing. Make it marshal the bookings that `f.mem.List` returns as one JSON array, and write them to `f.path` with `replaceFile`, returning every error. Your contract test runs against `FileStore` too.

> [!WARNING]
> A `save` that drops an error lets `Create` keep a booking its file never got, and the next restart loses it.
