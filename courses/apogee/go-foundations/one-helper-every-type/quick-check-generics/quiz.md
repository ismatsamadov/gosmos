## Why doesn't `Bigger` compile?

- [ ] Generic functions can't use operators like `>`
  > They can when the constraint allows it: with `cmp.Ordered` instead of `any`, `>` works.
- [x] `any` admits types with no order, so `>` isn't allowed on `T`
  > An operator works on `T` only if every type the constraint admits supports it, and `any` admits structs, maps and slices too.
- [ ] Every call must name the type, as in `Bigger[int]`
  > Go rejects the function itself, before any call, so a type argument wouldn't help.

## Does the call need a type argument, as in `Filter[Pass]`?

- [ ] Yes: every call to a generic function names its types
  > Only when Go can't work them out, and here the arguments tell it.
- [ ] No, because `Filter` turns every value into an `any`
  > Nothing is converted: `T` becomes `Pass`, and `high` is a `[]Pass`.
- [x] No: Go infers that `T` is `Pass` from `passes`
  > Go reads `T` off the arguments' types, so the call means `Filter[Pass](passes, aboveTen)`.
