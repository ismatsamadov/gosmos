## Which comparison of the two planners can go on the glass?

- [ ] The old planner timed on Monday and passplan 2 today, both on passplan's machine
  > Two moments differ in more than the planners: Monday's run may have shared the machine with other work, and nothing in its number says so.
- [ ] Each planner timed on its author's laptop, on the same windows
  > Two laptops differ in speed, and the gap between them would land in the ratio as if a planner had made it.
- [x] Both today, one after the other on passplan's machine, on the same windows, ten runs each
  > The same input, the same machine and the same hour leave the planners as the only difference, and ten runs show how far each number moves on its own.

## If passplan 2 planned 17 May in half the old planner's time, with 1,230 contacts to the old one's 1,231, what would Nkem write?

- [x] Nothing yet: a planner that gives a different answer isn't faster, it's wrong
  > First find the contact they disagree on and fix the planner that's wrong; only two planners that agree can be timed against each other.
- [ ] Twice as fast, one contact short
  > Half the time for a different plan compares two different jobs; the missing contact may be exactly the work that was skipped.
- [ ] Twice as fast: one contact in 1,231 is too few to matter
  > That one contact is a pass a satellite would lose, and a different answer means the two did different work.

## What does `BenchmarkMedian` measure from its second turn on?

- [ ] The same work in every turn, since `gusts` never changes its length
  > Its length stays, but `median` sorted its values in place on the first turn, so every later turn starts from sorted readings.
- [x] The median of readings already sorted: the best case
  > The first turn sorted `gusts` in place; a clone made with the timer stopped would give every turn the readings as recorded.
- [ ] Nothing: the compiler drops a call whose result isn't used
  > Inside `for b.Loop()`, every call and its result are kept, so `median` runs on every turn.
