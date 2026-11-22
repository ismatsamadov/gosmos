## What does `s` hold after `var s Station`?

- [ ] Nothing: it is `nil` until a value is assigned
  > A struct variable is never `nil`. Declaring it gives a usable value right away, with every field at its zero value.
- [x] `Code` is `""` and `MaskDeg` is `0`
  > A struct's zero value has each field at its own zero value.
- [ ] Whatever was left in that memory before
  > Go always starts a declared variable at its zero value, never at leftover data.

## After the last line, what is `a.MaskDeg`?

- [x] `5`
  > `b := a` copied every field, so changing `b.MaskDeg` leaves `a` alone.
- [ ] `9`
  > That would need `b` to refer to `a`, but assigning a struct makes a copy.
- [ ] `0`
  > `b := a` copies `a`; it doesn't move the value out of `a` or clear it.

## Which fields does `Station{Code: "PEL"}` set?

- [ ] None: it doesn't compile without `MaskDeg`
  > A literal with field names may leave fields out. Only a literal without names must list every field.
- [ ] `Code`, and `MaskDeg` has no value until one is assigned
  > Go has no unset fields: one that the literal leaves out holds its zero value, `0`.
- [x] `Code`, and `MaskDeg` holds `0`
  > A field that a literal leaves out holds its zero value.
