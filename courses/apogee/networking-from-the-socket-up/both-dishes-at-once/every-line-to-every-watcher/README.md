---
slug: every-line-to-every-watcher
title: Every Line to Every Watcher
type: code_tests
position: 3
difficulty: 8
topics: mutexes, networking
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Wednesday 26 April. Every second, dishd must send each dish's line to everyone on the watch port: Sunniva's terminal, Ingrid's screen and Tomás's laptop.

Watchers join and leave from goroutines of their own while lines go out, so the set lives behind a mutex. Hold the lock only long enough to copy the set, then write outside it: a slow write must never keep a watcher from joining. A failed write means that watcher hung up: remove it and go on.

```go
mu.Lock()
names := make([]string, 0, len(rota))
for name := range rota {
	names = append(names, name)
}
mu.Unlock()
// names is yours now: no lock needed
```

## Your task

`Broadcast` writes to nobody. Copy `w.conns` into a slice under `w.mu`, then write `line` and a newline to each watcher outside the lock, and `w.Remove` each one whose write fails.

> [!WARNING]
> Don't stop at the first failed write: every watcher after it would miss the line.
