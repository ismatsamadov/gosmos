## Why does a `min=0s` rule on `Turnaround` look at its type as well as its kind?

- [ ] It needn't: the kind `int64` already marks a duration
  > A plain `int64`, such as a count of frames, has the same kind, and "0s" is no number for it.
- [x] The kind says it holds an integer; only the type says to read "0s" as duration text
  > So a rule chooses by kind, then by type: the kind says to read the value with `Int`, and the type `time.Duration` says to read the bound with `time.ParseDuration`.
- [ ] The kind of `Turnaround` is `time.Duration`, the same as its type
  > A kind is what a type is built from: the fourth line prints `int64` for it.

## What does `f.Tag.Get("unit")` return, since `Turnaround` has no `unit` key?

- [x] "", which it would also return for a `unit` key that is there but empty
  > `Get` can't tell the two apart; `Lookup` can, with a second result that is `false` here.
- [ ] An error saying the key is missing
  > `Get` returns only a string, and `Lookup` adds a bool, never an error.
- [ ] The whole tag, since no key matched
  > `Get` returns the value of the one key it's asked for, and "" when it finds none.

## What would `t.FieldByName("Turnaround")` do if the second line left out `.Elem()`?

- [ ] Find the same field: Go follows the pointer for you
  > `TypeOf(&cfg)` is the type `*main.Config`, not `Config`, and `FieldByName` doesn't look through the pointer.
- [ ] Return `false` as its second result, with no field
  > It never gets that far: a pointer type has no fields to search.
- [x] Panic, because `*main.Config` is a pointer type and only `Elem()` reaches the struct
  > Its kind is `Pointer`, and `FieldByName` panics on any type that isn't a struct: "reflect: FieldByName of non-struct type *main.Config".
