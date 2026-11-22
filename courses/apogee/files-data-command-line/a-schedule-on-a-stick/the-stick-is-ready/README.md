---
slug: the-stick-is-ready
title: The Stick Is Ready
type: code_tests
position: 8
difficulty: 4
topics: files, testing
created: 2026-11-04T17:28:41Z
modified: 2026-11-04T17:28:41Z
---
Paz is ready to leave Port Ellery. Before the memory stick goes in the van, one function must prepare the report folder and save the complete day's schedule, with a useful error if either step fails.

Composition keeps each operation's job small: use the folder helper, construct the filename, then call the replacement helper. Return immediately when a prerequisite fails; a caller must never hear success after a failed setup.

```go
if err := prepareCache(root); err != nil {
    return err
}
return saveIndex(filepath.Join(root, "index.txt"))
```

## Your task

Complete `exportDay` using the supplied read-only `prepareReports` and `replaceSchedule` helpers. Save `text` at `root/reports/daily/<day>.txt`. Here `day` is a trusted label supplied by the program, not an arbitrary path from a user. Propagate setup and save errors.
