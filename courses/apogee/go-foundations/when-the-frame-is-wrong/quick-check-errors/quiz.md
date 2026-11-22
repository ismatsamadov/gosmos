## What must a type have to be used as an `error`?

- [ ] A field named `Message`
  > `error` asks for a method, not a field, and a field never satisfies an interface.
- [x] A method `Error() string`
  > `error` is an interface with exactly that one method, so any type that has it can be used as one.
- [ ] A method `String() string`
  > That's the method `fmt.Stringer` asks for; `error` asks for `Error() string`.

## Here `checkBattery` returns `nil`. What happens next?

- [x] The reading is stored
  > A `nil` error means nothing went wrong, so the `if` is skipped and `store` runs.
- [ ] `rejected:` prints with an empty reason
  > The `if` body runs only for an error that isn't `nil`.
- [ ] The program panics, because `err` holds nothing
  > Comparing `err` with `nil` is always safe; nothing is called on it.

## Which message follows Go's style for errors?

- [ ] `Battery too low`
  > Go error messages start in lowercase, because they're often printed after other text.
- [ ] `battery too low.`
  > Go error messages end without a period, because more text may follow them.
- [x] `battery too low`
  > Lowercase and no final period, so it reads well inside a longer line.
