---
slug: rotate-at-a-size
title: Rotate at a Size
type: code_tests
position: 2
difficulty: 8
topics: files, logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:38Z
---
At 08:25 on Monday you turned Lantern Atoll's uplinkd down to INFO on its debug port, and `uplinkd.log` all but stopped growing, at 41 GB: one file that nobody can open, let alone search. From Friday, uplinkd writes `uplinkd.log` only until the next record would take it past 64 MiB.

Then comes log rotation: the rotator closes the current file, renames it under its UTC day and that day's next number, `uplinkd-2028-12-04.1.log`, then `.2.log`, and opens a fresh `uplinkd.log`. A record is one `Write`, so a rotator that decides before each write never splits a line, and a record longer than the limit goes alone into a fresh file. A rename moves a file to a new name:

```go
old := filepath.Join(dir, "plan.txt")
err := os.Rename(old, old+".bak")
```

## Your task

1. Make `rotate` close `r.f`, rename the current file in `r.Dir` to `rotatedName(r.Program, r.day, r.n+1)`, and open a fresh one with `r.open()`, which also counts the day's files again.
2. Make `Write` call `rotate` first when the file isn't empty and `p` would take it past `r.MaxBytes`.
