## Does `earthRadiusKm + altKm` build?

- [x] Yes: the constant has no type of its own, so it works as a `float64`
  > An untyped constant takes the type the expression needs, here the `float64` of `altKm`.
- [ ] No: 6371 is an `int`, and Go won't add an `int` to a `float64`
  > That rule is for typed values, like an `int` variable; an untyped constant adapts.
- [ ] Yes, but the sum is rounded to a whole number
  > The constant works as a `float64`, so the sum is a `float64` and keeps any fraction.

## What happens with `const now = time.Now()`?

- [ ] `now` holds the moment the program started
  > A constant can't wait for the program to start: its value must be known when the program is compiled.
- [x] It doesn't compile: the time isn't known until the program runs
  > A constant's value is fixed when the program is compiled, and `time.Now()` only has an answer at run time.
- [ ] `now` gives the current time whenever it is read
  > Nothing can change a constant, so it can't follow the clock.
