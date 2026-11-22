---
slug: quick-check-the-plan
title: Quick Check: The Plan
type: quiz
position: 2
difficulty: 8
topics: context, goroutines, iterators
created: 2026-11-12T03:18:38Z
modified: 2026-11-12T03:18:38Z
---
Friday 22 June. Ingrid writes the checklist for Monday's rehearsal and, before she signs it, asks Nkem three questions about the planner in plain words.

Nkem answers with Tuesday's `Plan`, cut down to its shape:

```go
g, _ := newGroup(ctx, cfg.Slots)
for i, s := range n.Sites {
	g.Go(func(ctx context.Context) error {
		return plan(ctx, s, &lists[i])
	})
}
err := g.Wait() // a canceled run: ctx.Err()
merged := MergeBy(byContact, seqs...)
```

## Your task

Answer Ingrid's three questions: why twelve stations can be planned at once, how many contacts the last merge holds at a time, and what the test Board shows when a run is canceled.
