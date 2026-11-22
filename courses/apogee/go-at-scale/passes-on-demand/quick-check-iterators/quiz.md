## When does the call `yield(3)` return?

- [ ] At once: the 3 waits for the loop, as a value in a channel would
  > Nothing waits anywhere: `yield` is the loop's body, made into a function, and it runs the body before it returns.
- [x] After the loop's body has printed 3
  > The body runs inside the call to `yield`, so the call returns when the body ends, with true if the loop goes on.
- [ ] After the loop has printed 3, 2 and 1
  > `countdown` calls `yield` once per value, and each call returns before the next value is made.

## A loop prints 3, then breaks when `n` is 2. What happens?

- [x] `yield(2)` returns false, and `countdown` returns without calling it again
  > `break` ends the body and makes its `yield` return false; a right iterator stops there, as `countdown` does.
- [ ] `yield(2)` returns true, and Go stops calling `countdown`
  > Go can't stop a function from outside: the false `yield` returns is how the iterator learns the loop is over.
- [ ] `countdown` goes on to `yield(1)`, and Go throws the value away
  > An iterator that called `yield` again would make Go panic; `countdown` checks the result and returns.

## After a loop has printed 3, 2 and 1, a second loop ranges over `countdown`. What is printed next?

- [ ] Nothing more: the first loop used the sequence up
  > That happens to an iterator that reads something that runs out, such as a reader; `countdown` makes its values afresh in each call.
- [ ] `go` and nothing else, since `n` is already 0
  > `n` belongs to each call of `countdown`, and every call starts it at 3.
- [x] `go`, then 3, 2 and 1 again
  > Each loop calls `countdown` from the beginning, so the second prints exactly what the first did.
