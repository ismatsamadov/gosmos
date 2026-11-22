`trace.NewTask` returns two things, the new context and the task, so assign both at the top of `runPass`. Naming the context `ctx` again means every later use of `ctx`, the `select` and `fly` included, carries the task with no other change.

Right after it, defer the task's `End` method, so the task ends however `runPass` returns. Then call `trace.Log` with the task's context, the category `"pass"` and `p.ID`.

If the test says `fly` doesn't get the task's context, the call to `u.fly` still names the context `runPass` was given, not the one `NewTask` returned. Remove the `_ = trace.NewTask` line once the task is there.
