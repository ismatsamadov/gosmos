---
slug: never-wait-forever
title: Never Wait Forever
type: code_tests
position: 0
difficulty: 7
topics: networking, time
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Ingrid's log of the storm shows dishd asking Dish 1 for its status at 03:10, then nothing until 03:50. A read with no deadline waits as long as the connection lives, and the cut fiber had left it dead but open.

`SetReadDeadline(t)` sets an instant after which a read on the connection gives up with an error matching `os.ErrDeadlineExceeded`. Set it to now plus your limit just before each read: one set when the connection opened soon passes while it's still in use.

```go
a, _ := net.Pipe() // nothing writes the other end
a.SetReadDeadline(time.Now().Add(time.Second))
_, err := a.Read(make([]byte, 16))
fmt.Println(err) // read pipe: i/o timeout
```

## Your task

`recv` reads a frame with no deadline at all. Replace its `TODO` so that the read gives up at now plus `c.Timeout`, 5 seconds.

> [!NOTE]
> The tests run inside `testing/synctest`, a bubble with a fake clock: every deadline keeps its real length and passes in no real time.
