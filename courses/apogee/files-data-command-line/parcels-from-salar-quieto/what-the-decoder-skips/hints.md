By the time your checks run, decoding has succeeded, so each check looks at one field of `m`: a missing element left that field at its zero value, with no error to tell you.

`Dish` is empty when `<dish>` is missing or empty, `Made.IsZero()` is true when `<made>` is missing, and `Files` has no entries when no `<file>` was found.

Return at the first gap with `errors.New` and a message that names the element, and return `m` with a nil error only after all three checks pass.
