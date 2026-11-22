---
slug: settings-from-a-file
title: Settings from a File
type: code_tests
position: 0
difficulty: 7
topics: errors, files, structs
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T17:14:45Z
---
The night shift wants its own prompt, and Port Ellery's passes when they type `passes` alone, set in a file they can change themselves:

```json
{
  "prompt": "night> ",
  "cache_ttl": "5m",
  "station": "PEL"
}
```

Decoding into a struct that already holds the defaults changes only the fields the file names, so every field is optional. In `settings.go`, `decodeStrict` refuses unknown fields and anything after the object, and `ttlText` refuses bad durations.

## Your task

1. Write `loadSettings`: a missing file gives the defaults; otherwise decode the file over them with `decodeStrict`, and wrap a decoding error with the path.
2. Write `apply`, which sets the shell's prompt, its cache's `ttl` and a new `station` field on `Shell`.
3. Make `passes` with no code use `s.station` as if it were typed; empty, it still means every station.
