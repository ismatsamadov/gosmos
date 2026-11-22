## The pager must stay quiet for `dup`. Which check recognizes it?

- [ ] `dup == ErrDuplicate`
  > `==` compares only the outer error, the wrapper made by `fmt.Errorf`, so it's `false`.
- [x] `errors.Is(dup, ErrDuplicate)`
  > `errors.Is` looks through every layer of wrapping for that exact value.
- [ ] `dup.Error() == "duplicate frame"`
  > The wrapper's message starts with the satellite and sequence number, so the texts differ.

## The pager needs the `Value` of the `*LimitError` inside `hot`. Which expression gets it?

- [ ] `errors.Is(hot, &LimitError{})`
  > `errors.Is` compares with one particular value, and a new `&LimitError{}` is a different pointer; it also returns only a `bool`.
- [ ] `hot.(*LimitError)`
  > A type assertion looks only at the outer error, the `fmt.Errorf` wrapper, so this one panics.
- [x] `errors.AsType[*LimitError](hot)`
  > It searches every layer of wrapping for a `*LimitError` and returns it, fields and all.
