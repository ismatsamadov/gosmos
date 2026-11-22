## Does snippet A have a data race?

- [ ] No: `screens++` is a single statement
  > One statement is still a read, an add and a write, and the two goroutines' steps can interleave.
- [x] Yes: two goroutines touch `screens`, both write, and nothing orders them
  > That is a data race: two goroutines, one variable, at least one writing, and no lock, channel or atomic operation between them.
- [ ] Only on a machine with more than one core
  > The race is in the program, not the machine: one core can switch goroutines between the read and the write, and `go run -race` reports it there too.

## Why does snippet B end with `count` at 2 once both have run?

- [x] Atomic operations order each other, so the program behaves as if one `Add` ran after the other
  > A program without races behaves as if its goroutines' steps were taken one at a time in some order: locks, channel operations and atomics give that order.
- [ ] The goroutine always runs before the line after it
  > Either can run first; the result is the same in both orders.
- [ ] `count` is a variable of the whole program
  > Where a variable lives changes nothing; the atomic operations do.

## What may a program with a data race do?

- [ ] Lose a count now and then, but nothing worse
  > A lost count is only the mildest outcome: a race on a string, a slice or a map can break far more.
- [ ] Run correctly whenever the race detector is off
  > The race detector only reports races; turning it off leaves them where they are.
- [x] Anything: a reader can even see half of a string another goroutine is writing
  > Go promises almost nothing for a racy program: a string or a slice can be read half old and half new, and that can corrupt memory.
