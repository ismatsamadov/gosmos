Test for a duration first, comparing `v.Type()` with `reflect.TypeFor[time.Duration]()`, then use `v.CanInt()` and `v.CanFloat()`; the stub's `return` stays at the end for every other type.

`strconv.ParseInt(text, 10, 64)` reads an int's bound as an `int64` to compare with `v.Int()`. For a duration, convert `v.Int()` with `time.Duration(...)` before comparing and returning it. An error from a parser means the bound doesn't fit.

If a duration's problem says the rule does not fit `time.Duration`, the int case caught it first: a `time.Duration` passes `CanInt` too, so its case must come before the int case.
