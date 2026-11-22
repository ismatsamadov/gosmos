## What does `fmt.Println(len(frames))` print?

- [ ] `1`
  > 1 is the index of the last element; `len` counts the elements.
- [x] `2`
  > The slice holds two elements, 4 and 8, so its length is 2.
- [ ] `8`
  > 8 is the value of the last element, not how many elements there are.

## What happens at `codes[5]`?

- [ ] The build fails: index 5 is out of range
  > The compiler doesn't check slice indexes against the length. The check happens when the line runs.
- [ ] It gives `""`, the zero value of a string
  > Go never reads past the end of a slice. It stops the program instead.
- [x] The program panics when the line runs
  > `codes` has indexes 0 to 2, so index 5 stops the program with an index out of range panic.

## Why write `s = append(s, x)` and not just `append(s, x)`?

- [x] `append` returns the grown slice, and `s` keeps its old length
  > Only the assignment puts the slice with `x` added into `s`.
- [ ] The assignment makes a safe copy of `s` first
  > The assignment copies nothing for safety. It is there to keep the slice `append` returns.
- [ ] Style only: both forms grow `s`
  > On its own, `append(s, x)` doesn't even compile, because its result must be used.
