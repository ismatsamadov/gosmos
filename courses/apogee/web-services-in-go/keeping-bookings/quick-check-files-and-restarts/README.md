---
slug: quick-check-files-and-restarts
title: Quick Check: Files and Restarts
type: quiz
position: 6
difficulty: 8
topics: errors, files
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Friday's deploy is the last before the holidays, and the backend desks stay empty for two weeks after it. Nkem's last three questions are about what `bookings.json` goes through, starting with how `replaceFile` saves it:

```go
tmp, err := os.CreateTemp(dir, ".bookings-*")
// write the bookings to tmp, close it, then:
err = os.Rename(tmp.Name(), path)
```

## Your task

Answer Nkem's three questions: what a reader of the file can see during a save, what a broken file must do to a start, and which ID comes first after a restart.
