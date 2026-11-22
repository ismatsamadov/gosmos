---
slug: writes-wait-too
title: Writes Wait Too
type: code_tests
position: 4
difficulty: 8
topics: networking, time
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Thursday 13 April. In Tuesday's trace, Nkem finds dishd stuck for six minutes in the `Write` of a POINT for Dish 1, its bytes going nowhere. "Reads have deadlines now. Do writes?"

A `Write` returns once the connection has taken your bytes, and a connection whose other end has stopped reading stops taking them once its buffers are full. `SetWriteDeadline` bounds that wait as `SetReadDeadline` bounds a read; each covers its own direction only:

```go
a, _ := net.Pipe() // nothing reads the other end
a.SetWriteDeadline(time.Now().Add(time.Second))
_, err := a.Write([]byte("STOW"))
fmt.Println(err) // write pipe: i/o timeout
```

## Your task

`send` writes with no deadline. Before it writes, set the connection's write deadline to now plus `c.Timeout`, so a controller that stops reading costs dishd 5 seconds and an error.

> [!WARNING]
> A read deadline doesn't help here: it bounds reads only. `SetDeadline` sets both at once.
