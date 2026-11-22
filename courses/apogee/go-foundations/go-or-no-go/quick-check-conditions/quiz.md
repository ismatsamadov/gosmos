## Which of these `if` statements compiles?

- [ ] `if windKmh >= 72: fmt.Println("stow")`
  > Go has no colon form: the statements an `if` guards always go between braces.
- [ ] `if (windKmh >= 72) fmt.Println("stow")`
  > Parentheses around the condition are allowed but unneeded, and the braces are still required.
- [x] `if windKmh >= 72 { fmt.Println("stow") }`
  > The statement sits between braces, which Go requires even around one statement; writing the whole `if` on one line is allowed.

## Which value of `maxElevDeg` makes `maxElevDeg > 5 && maxElevDeg < 10` true?

- [ ] `5`
  > `>` leaves out the value itself, so 5 > 5 is false, and so is the whole condition.
- [x] `7`
  > 7 > 5 and 7 < 10 are both true, and `&&` needs both.
- [ ] `12`
  > 12 > 5 is true, but 12 < 10 is false; `&&` needs both sides true, unlike `||`.

## With `up` false, what does `up && dishReady()` do?

- [x] It is false, and `dishReady` never runs
  > `&&` stops at a false left side, since the result can no longer be true, so nothing prints.
- [ ] It prints `asking the dish`, then is false
  > Go skips the right side of `&&` when the left side has already decided the result.
- [ ] It prints `asking the dish`, then is true
  > `&&` is true only when both sides are, and the left side is false; the call is skipped as well.
