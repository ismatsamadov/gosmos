Only the type of `n` changes, to `atomic.Uint64`. Its zero value is a count of 0, so a new `Counter` needs no setup, and the `_ = new(atomic.Uint64)` line can go once `Inc` uses the field.

`Inc` and `Add` call the field's `Add` method, with 1 and with `n`, and `Value` returns what its `Load` method reads: each method is one line.

If a count comes up short, `Inc` still takes more than one step: `n++` on a plain `uint64` does, and so does a `Load` followed by a `Store`. If the counts are right but the test still names `metrics.go`, `n` is still a plain `uint64`: a test can't catch a lost count every time, but it can read the field's type.
