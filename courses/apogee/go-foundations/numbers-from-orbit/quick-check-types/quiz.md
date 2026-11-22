## What type does `beaconMHz` get?

- [ ] `int`
  > An `int` holds whole numbers only, and 437.25 has a fraction.
- [x] `float64`
  > `:=` takes the type from the value, and a number written with a decimal point gets `float64`.
- [ ] `string`
  > It would be a `string` only in quotes, as `"437.25"`.

## What does `operator` hold before anything is assigned to it?

- [x] `""`, the empty string
  > A variable declared without a value holds its type's zero value, and a `string`'s is `""`.
- [ ] `0`
  > `0` is the zero value of numbers; a `string` starts out as empty text.
- [ ] No value yet: reading it is a compile error
  > Go gives a variable declared without a value its type's zero value, so reading it before any assignment is fine.

## What does `:=` do?

- [ ] Gives a new value to a variable that already exists
  > That is what `=` does. After `var operator string`, writing `operator := "Apogee"` doesn't compile, since `operator` already exists.
- [ ] Declares a variable and leaves it at its zero value
  > That is `var` with a type and no value, like `var operator string`.
- [x] Declares a new variable and gives it a value
  > `:=` creates the name and sets its first value, taking the type from that value.
