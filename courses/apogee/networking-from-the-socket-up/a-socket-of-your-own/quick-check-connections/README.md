---
slug: quick-check-connections
title: Quick Check: Connections
type: quiz
position: 2
difficulty: 6
topics: defer, interfaces, networking
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T07:44:27Z
---
Tuesday. Before dishd talks to anything real, Nkem sends three questions from headquarters at Port Ellery, about the start of a function that reads from a dish controller:

```go
conn, err := net.Dial("tcp", addr)
if err != nil {
	return err
}
defer conn.Close()
buf := make([]byte, 64)
n, err := conn.Read(buf)
```

## Your task

Answer Nkem's three questions: what `conn` is besides a connection, what `n` can be once the controller has sent its banner, and why the `defer` stands where it does.
