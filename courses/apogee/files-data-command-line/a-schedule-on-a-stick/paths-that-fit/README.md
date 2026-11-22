---
slug: paths-that-fit
title: Paths That Fit
type: code
position: 3
difficulty: 3
topics: files
created: 2026-11-04T17:28:41Z
modified: 2026-11-07T19:09:32Z
---
Paz wants one folder per day on the stick. Hand-written slashes have already left a double separator in one path and a leftover `..` in another.

`filepath.Join` combines local path components and cleans repeated separators, `.` and parent steps such as `..`. It builds a path string; it does not create a folder or check that anything exists:

```go
path := filepath.Join(
    "logs", "old", "..", "current.txt")
fmt.Println(path) // logs/current.txt on this Mac
```

## Your task

Replace the two placeholders with `filepath.Join` calls. First join `plans`, `2027-07-12` and `passes.txt`. Then join `plans`, `old`, `..` and `report.txt`. Keep one printed path per line and remove the import placeholder.
