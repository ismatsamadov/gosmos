Copying the `Schedule` value copies its fields, not its contacts, so the copy's window can change without touching `old`. Set the replan's `Version` to 1 and `Made` to `now()`, and `Offer` it to a new `Publisher` before `newReplanner` reads it.

Keep the last pass flown, starting from a `PassView` whose AOS is `at` and whose ID is empty, so the old planner's contacts count as flown. The next contact is the first that `byAOS` puts after that pass; a `time.After` for the time from `now()` to its AOS wakes then, and a nil channel never does, once none is left.

Select on that timer, `r.news` and `ctx.Done()`, reading `pub.Current()` again each time round. If the 25-minute test finds a contact on a dish the old planner holds, run 1 was planned from the hand-over alone and later replans forgot the old planner's contacts.
