## What does this code print?

- [ ] `118`
  > `nextSeq` adds one to its own copy of `seq`; the caller's variable never changes.
- [x] `117`
  > Arguments are copies: `nextSeq` changes its parameter, and the value it returns is thrown away.
- [ ] Nothing: the build fails, because the result of `nextSeq` is never used
  > Go lets a call ignore a function's results. An unused variable stops the build, but an ignored result doesn't.

## What must `f` return for `a, b := f()` to compile?

- [x] Exactly two values
  > Each name on the left takes one result, so the counts must match.
- [ ] Two or more values, and the extra ones are dropped
  > An assignment never drops results on its own; to skip one, receive it with `_`.
- [ ] One value, and `b` gets its zero value
  > A function with one result can't fill two names, so the build fails.

## Which line declares a function that takes an `int` and returns a `string`?

- [ ] `func string label(n int)`
  > That puts the result type before the name, as C does; Go puts it after the parameters.
- [ ] `func label(n int): string`
  > Go writes no colon before the result type.
- [x] `func label(n int) string`
  > The result type comes right after the parameter list.
