---
slug: make-room-for-reports
title: Make Room for Reports
type: code_tests
position: 4
difficulty: 3
topics: files, testing
created: 2026-11-04T17:28:41Z
modified: 2026-11-07T19:09:32Z
---
The van's laptop has no reports folder on its first run. Paz needs the tool to prepare it, and to keep working when the same folder is already there tomorrow.

`os.MkdirAll` creates a directory and any missing parents. An existing directory is fine; a regular file blocking the path is an error. The mode supplies creation permissions, subject to the operating system's permission mask:

```go
dir := filepath.Join("cache", "images")
err := os.MkdirAll(dir, 0700)
if err != nil {
    return err
}
```

## Your task

Implement `prepareReports`. Join `root`, `reports` and `daily`, then create that path with mode `0700`. Return the error so the caller cannot mistake a blocked path for a ready folder.
