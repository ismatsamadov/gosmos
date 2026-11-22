---
slug: read-the-schedule
title: Read the Schedule
type: code_tests
position: 0
difficulty: 3
topics: files, testing
created: 2026-11-04T17:28:41Z
modified: 2026-11-04T17:28:41Z
---
Paz hands you a memory stick with the van's next passes. Beyond Port Ellery there will be no signal, so the schedule must come from the file, not a service the laptop might not reach.

`os.ReadFile(path)` reads the whole file and returns its bytes and an error. Check the error before using the bytes. `string(data)` turns a byte slice into text without trimming spaces or line endings:

```go
data, err := os.ReadFile("notes.txt")
if err != nil {
    return err
}
fmt.Println(len(data))
```

## Your task

Implement `readSchedule`. Return the file's exact text on success, or an empty string and the original error when reading fails. Remove the import placeholder. Tests create their own temporary files; no real memory stick is needed.
