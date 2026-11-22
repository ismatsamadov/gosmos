---
slug: save-the-report
title: Save the Report
type: code_tests
position: 1
difficulty: 3
topics: files, testing
created: 2026-11-04T17:28:41Z
modified: 2026-11-07T19:09:32Z
---
Ingrid needs a report the van can bring back even if the road stays offline. An empty file must mean an empty report, not a save that silently did nothing.

`os.WriteFile` writes bytes and returns an error. It creates a missing file or replaces the contents of an existing one. Its permission argument applies when creating a file; it does not change an existing file's permissions:

```go
data := []byte("Ready\n")
err := os.WriteFile("note.txt", data, 0600)
if err != nil {
    return err
}
```

`0600` permits the owner to read and write, with no access for others.

## Your task

Implement `saveReport` using `os.WriteFile`, converting `text` to bytes and using creation mode `0600`. Return its error. The parent folder must already exist; a later lesson creates folders.
