---
slug: a-slow-watcher-waits-alone
title: A Slow Watcher Waits Alone
type: code_tests
position: 4
difficulty: 9
topics: mutexes, networking, time
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Thursday 27 April. Tomás's laptop went to sleep with the watch port open, and every other screen froze with it: `Broadcast` sat in its write to the laptop, and no line went anywhere else.

A connection nobody reads takes bytes only until its buffers fill; then a write to it waits, and in a loop over every watcher, everyone after it waits too. A write deadline turns that wait into an error, and the watcher behind it is removed while the loop goes on:

```go
c.SetWriteDeadline(time.Now().Add(time.Second))
_, err := io.WriteString(c, "YST-1 ...\n")
late := errors.Is(err, os.ErrDeadlineExceeded)
// late: the write gave up after a second
```

## Your task

`Broadcast` writes with no deadline. Before each write, set that watcher's write deadline to now plus `w.Timeout`. A failed write still removes its watcher.

> [!WARNING]
> Give each write its own deadline, set just before it: one deadline for the whole loop is used up by the first slow watcher.
