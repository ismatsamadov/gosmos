## Once `lines` is built, before the loop starts, how many contacts have been read?

- [x] None: building the chain only wraps each function in the next
  > Each adapter returns a function and runs nothing; the first contact is read when the loop ranges over `lines`.
- [ ] All 260, filtered and kept for the loop
  > Nothing is kept anywhere: `Filter` yields as it reads, and reads only once something ranges over it.
- [ ] Three, the ones `Take` will yield
  > `Take` can't know which contacts will reach it until the loop asks, so building it reads none.

## In what order does the work happen once the loop runs?

- [ ] `Filter` reads all 260, then `Map` makes the lines, then `Take` keeps three
  > That would be stage by stage; an adapter hands on one value at a time, inside the call to its `yield`.
- [x] One contact at a time: each is read and tested, and a kept one is made a line, taken and printed, before the next is read
  > Each call to `yield` runs the next adapter's body, so a contact travels the whole chain before the schedule makes another.
- [ ] `Take` asks for three lines, then `Map` and `Filter` make them in one go
  > Nothing is asked for ahead: each value is handed up the chain as soon as the schedule makes it.

## What happens below `Take` once it has yielded its third line?

- [ ] `Filter` and the schedule run to the end of the day, and their values are thrown away
  > Then a short chain would cost as much as a long one: every adapter returns when its `yield` returns false.
- [ ] Go ends the chain from outside, since the loop has its three lines
  > Go can't stop a function from outside; each adapter returns because its `yield` returned false.
- [x] `Take` returns, so `Map`'s `yield` returns false and `Map` returns, then `Filter`, and the schedule's `yield` returns false
  > Each adapter's loop ends when the one above it returns, so the stop travels down to the source.
