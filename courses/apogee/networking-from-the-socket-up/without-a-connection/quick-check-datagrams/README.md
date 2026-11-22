---
slug: quick-check-datagrams
title: Quick Check: Datagrams
type: quiz
position: 3
difficulty: 8
topics: networking
created: 2026-11-13T06:15:46Z
modified: 2026-11-13T06:15:46Z
---
Tuesday afternoon. Nkem has read the live view and the recorder's packing from Port Ellery, and sends her questions one at a time, as ever. They're about this exchange between two sockets with nothing else on the line:

```go
rec.WriteTo(a, view.LocalAddr()) // 72 bytes
rec.WriteTo(b, view.LocalAddr()) // 70 bytes
buf := make([]byte, 64)
n1, _, _ := view.ReadFrom(buf)
n2, _, _ := view.ReadFrom(buf)
```

## Your task

Answer Nkem's three questions: what the second read returns, what a `WriteTo` without an error promises, and why the recorder packs whole lines.
