---
slug: the-offline-resource-pack
title: The Offline Resource Pack
type: code_tests
position: 8
difficulty: 8
topics: errors, files, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-22T10:21:59Z
---
On Monday the van reaches Salar Quieto. Paz's stick holds `calibration.json`: how far the van's own dish pointed off when it was checked at Port Ellery before the trip. gctl's default is zero. If the stick's file is broken, quietly using zero would copy the van dish's error into La Grande.

An override falls back to the default only when it has no file, which `errors.Is(err, fs.ErrNotExist)` tells. Any other error, a file that does not parse among them, must reach the caller:

```go
data, err := fs.ReadFile(stick, "motd.txt")
if errors.Is(err, fs.ErrNotExist) {
	data, err = fs.ReadFile(builtIn, "motd.txt")
}
```

## Your task

Complete `loadResources`: read `calibration.json` from `override`, or from `fsys` when `override` has none, parse it with `parseOffsets`, and store the result in `Reference`. Return every other error.
