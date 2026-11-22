`now` is a function: calling it, `now()`, returns a `time.Time` like `time.Now()` does, so the rest of `checkAt` stays as it is.

Only one field of the options changes: `CurrentTime` gets the result of calling `now`, parentheses included.

If the test still fails at `NotBefore` exactly, `CurrentTime` still reads the Mac's clock: look for any `time.Now()` left in `checkAt`.
