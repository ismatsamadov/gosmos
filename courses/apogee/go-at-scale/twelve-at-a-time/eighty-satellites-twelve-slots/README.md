---
slug: eighty-satellites-twelve-slots
title: Eighty Satellites, Twelve Slots
type: code_tests
position: 7
difficulty: 9
topics: context, errors, goroutines, mutexes
created: 2026-11-11T12:51:42Z
modified: 2026-11-11T12:51:42Z
---
Friday 16 March. For the first time, passplan 2's predict stage runs on its new machine: eighty satellites over twenty stations, for take-over day.

On Monday `predictAll` took a slot for each satellite, but a failed prediction stopped nothing. Thursday's group stops at its first failure. Run each satellite's prediction as a job of a group, and collect the passes in a map guarded by a `sync.Mutex`, since jobs can finish at the same moment.

```go
g, _ := newGroup(ctx, 2)
for _, code := range []string{"YST", "TIR"} {
	g.Go(func(ctx context.Context) error {
		return stow(ctx, code)
	})
}
err := g.Wait() // the first stow that failed
```

## Your task

Rewrite `predictAll` on the given `Group`, with `slots` as its limit: one job for each satellite, which predicts it and adds its passes to the map under a mutex. Return the map, or `nil` and the error `Wait` returns.

> [!NOTE]
> The test also runs it for real: take-over day's 7,485 passes of all eighty satellites over twenty stations, Furrow-A's 105 among them.
