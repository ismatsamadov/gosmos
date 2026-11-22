---
slug: quick-check-caches
title: Quick Check: Caches
type: quiz
position: 4
difficulty: 7
topics: mutexes, time
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T17:14:45Z
---
Before the cache meets a real night, Nkem asks three questions about the lines that decide its age and what it hands out:

```go
fresh := c.now().Before(c.loadedAt.Add(c.ttl))
if !c.loaded || !fresh {
	if err := c.reload(); err != nil {
		return nil, err
	}
}
return slices.Clone(c.cached), nil
```

## Your task

Answer three questions: when passes loaded at 04:00 go stale, why the clock is a field, and why `passes` returns a copy.
