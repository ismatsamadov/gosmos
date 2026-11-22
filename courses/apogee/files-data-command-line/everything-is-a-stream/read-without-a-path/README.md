---
slug: read-without-a-path
title: Read Without a Path
type: code_tests
position: 0
difficulty: 3
topics: interfaces, streams, testing
created: 2026-11-04T17:36:06Z
modified: 2026-11-05T08:00:06Z
---
The van's recorder can supply a file or a small test recording. Paz wants the same text reader to accept either, without knowing a filename.

`io.Reader` describes something that supplies bytes. Files, `strings.Reader` and many other types satisfy it. `io.ReadAll` consumes a reader and returns its bytes and error. A reader can produce some bytes before failing; keep the error even when there is text:

```go
r := strings.NewReader("PEL ready\n")
data, err := io.ReadAll(r)
fmt.Println(len(data), err)
```

## Your task

Implement `readText` with `io.ReadAll`. Return the bytes as text together with the read error. Partial text may help diagnose a failed recording, but a non-nil error still means the read failed.
