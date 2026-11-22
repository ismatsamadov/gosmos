Put a `switch` on `rule` with a `case "required":` above the `return` that reports an unknown rule, so every other rule still reaches it.

`IsZero` covers the blank name, the 0 dishes and a nil list. An emptied list isn't zero, so when the value's `Kind()` is a slice, its `Len()` decides too.

If a test reports a panic from `Len` on an int, call `Len` only for a slice: it panics on a number.
