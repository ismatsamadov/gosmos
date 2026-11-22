---
slug: quick-check-byte-streams
title: Quick Check: Byte Streams
type: quiz
position: 3
difficulty: 6
topics: networking, streams
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Tuesday evening. Before dishd reads its first real frame, Nkem sends three questions from Port Ellery about what a connection carries. Two frames are on their way, a PING and a QUERY that asks for the dish's status, 11 bytes each:

```go
conn.Write(ping)
conn.Write(query)
```

## Your task

Answer Nkem's three questions: how the 22 bytes can reach the controller, what `io.ReadFull` does after a short read, and why every frame carries its length when the banner ends with a newline.
