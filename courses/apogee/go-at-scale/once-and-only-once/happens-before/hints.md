The planner's goroutine raises the flag, waits for the Board, and only then sets `plan.Version` and `plan.Contacts`. The Board reads the plan as soon as it sees the flag.

Cut the two lines that set `plan.Version` and `plan.Contacts` and paste them right after `defer close(done)`, before `ready.Store(true)`.

If the output still says `plan v0: 0 contacts`, the plan is filled after the `Store`, or after the wait for `looked`, which only ends once the Board has printed.
