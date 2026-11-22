## How many times does the loop above print `ping`?

- [ ] 2
  > `i` starts at 0, not 1, and `i < 3` is still true when `i` is 2, so the loop prints for 0, 1 and 2.
- [x] 3
  > `i` takes the values 0, 1 and 2; when it reaches 3, `i < 3` is false and the loop ends.
- [ ] 4
  > The condition is checked before each round, so the loop ends when `i` reaches 3, without printing; `i <= 3` would print four times.

## What does `for { fmt.Println("ping") }` do?

- [x] It prints `ping` forever
  > With no condition, nothing ends the loop: only a `break` or a `return` inside it could.
- [ ] It prints `ping` once
  > Running a block once is what an `if` does; a `for` repeats its block.
- [ ] It doesn't compile without a condition
  > Go accepts a `for` with no condition at all; it repeats until something inside it leaves the loop.
