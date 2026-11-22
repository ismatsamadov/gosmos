---
slug: append-never-truncate
title: Append, Never Truncate
type: code_tests
position: 0
difficulty: 7
topics: files, logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-12T11:41:59Z
---
The page reads `LAN: disk 98% full`, and the 41 GB file is `uplinkd.log`. The new log files start with how a file is opened: a restart in the middle of a day must add to the day's file, never empty it.

`os.Create` empties a file that already exists. `os.OpenFile` takes flags instead, joined with `|`: with `O_APPEND` every write goes to the file's end, and with `O_CREATE` a missing file is made with the mode you give. A log file's mode is `0640`: its owner reads and writes it, its group only reads it, and other users get no access.

```go
// make a lock file, unless one is there
f, err := os.OpenFile("passplan.lock",
	os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
if errors.Is(err, fs.ErrExist) {
	fmt.Println("passplan is already running")
}
```

## Your task

Make `openLog` open the path it's given with `os.OpenFile`, write-only and appending, creating the file with mode `0640` when it is missing.

> [!NOTE]
> The mode counts only when the file is created, minus any bits the system's umask setting holds back; a file that exists keeps its own.
