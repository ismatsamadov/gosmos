## Which goroutines does the pass's task cover?

- [ ] Only the goroutine that called NewTask
  > The task travels in `ctx`, so what the new goroutine does with that context belongs to the task as well.
- [x] Both, since the new goroutine runs with the task's context
  > A task can span goroutines: the region and `fly`, run with its context in the new goroutine, count as the task's.
- [ ] Only the new goroutine, where the region is
  > The task starts with `NewTask` in the first goroutine, and runs from there until `task.End()`.

## Where must r.End() run?

- [x] In the goroutine that started the region, once any region started inside it has ended
  > A region stays in one goroutine and regions nest; the deferred `r.End()` in the same function keeps both rules.
- [ ] In any goroutine, as long as it has the task's context
  > The context carries the task, not the region: a region belongs to the goroutine that started it.
- [ ] In the goroutine that called NewTask, when the task ends
  > That goroutine never started the region, and a region must end in the goroutine that started it.

## With no trace running, what do the task and the region cost?

- [ ] As much as with a trace, since tasks are always recorded
  > Nothing is recorded when no trace runs: `trace.IsEnabled()` is false, and the calls return almost at once.
- [ ] Nothing at all, since the compiler removes calls to trace
  > The calls still run and do almost nothing; the compiler leaves them in place.
- [x] Almost nothing: tasks and regions record nothing while no trace runs
  > With `trace.IsEnabled()` false they skip their work, which is why uplinkd can keep them in every pass.
