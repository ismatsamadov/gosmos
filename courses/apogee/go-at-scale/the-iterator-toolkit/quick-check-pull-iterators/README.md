---
slug: quick-check-pull-iterators
title: Quick Check: Pull Iterators
type: quiz
position: 7
difficulty: 8
topics: iterators
created: 2026-11-11T22:52:57Z
modified: 2026-11-11T22:52:57Z
---
Nkem reads your merge with a finger on each head, then asks three questions about the pull iterators under it.

Her whiteboard loop pulls the pass IDs of one station's day, one at a time, until `next` says the sequence has ended:

```go
next, stop := iter.Pull(slices.Values(ids))
defer stop()
for {
	id, ok := next()
	if !ok {
		break
	}
	fmt.Println(id)
}
```

## Your task

Answer Nkem's three questions: what `next` gives after the end, when `stop` must run, and who may call the two.
