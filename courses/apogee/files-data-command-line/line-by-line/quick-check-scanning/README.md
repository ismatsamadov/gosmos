---
slug: quick-check-scanning
title: Quick Check: Scanning
type: quiz
position: 3
difficulty: 4
topics: errors, streams
created: 2026-11-06T16:34:04Z
modified: 2026-11-06T16:34:04Z
---
The van finds a signal at a bend in the road, and Nkem reviews the scanner before it meets all 200,000 lines. One draft on the laptop builds, but has two mistakes:

```go
sc := bufio.NewScanner(log)
n := 0
for sc.Scan() {
    n++
}
sc.Buffer(make([]byte, 0, 4096), 256*1024)
return n, nil
```

## Your task

Answer Nkem's three questions about this draft: what the end of its loop proves on its own, where `Buffer` belongs, and what a record past the limit does.
