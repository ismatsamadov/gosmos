---
slug: quick-check-sharing
title: Quick Check: Sharing
type: quiz
position: 6
difficulty: 4
topics: slices
created: 2026-11-01T07:36:45Z
modified: 2026-11-22T10:08:26Z
---
Ingrid, at the Board: "Walk me through how a pass moves when nobody moves it." She writes five lines on the glass:

```go
a := []string{"PEL", "SQT", "LAN"}
b := a[:2]
b[0] = "X"
c := slices.Clone(a)
c[1] = "Y"
```

Two questions, both about `a`.

## Your task

Follow the lines in order, and for each write, ask which array it lands in before you answer.
