## What is `a[0]` after `b[0] = "X"`?

- [ ] `"PEL"`
  > That would be true if `a[:2]` copied, but `b` shares `a`'s array, so the write shows in `a`.
- [x] `"X"`
  > `b` is a view of the same array as `a`, so writing `b[0]` writes `a[0]`.

## What is `a[1]` after the last line?

- [x] `"SQT"`
  > `slices.Clone` gave `c` an array of its own, so `c[1] = "Y"` changes only `c`.
- [ ] `"Y"`
  > That would need `c` to share `a`'s array, but `slices.Clone` copies the elements.
