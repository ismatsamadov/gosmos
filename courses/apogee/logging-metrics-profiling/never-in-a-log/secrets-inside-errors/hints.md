Declare a variable at the top level of `redact.go`, outside every function, and set it with `regexp.MustCompile` from `secretPattern`, as the example sets `orb`: Go sets it once, when the program starts, and `scrub` can use it on every call. Then remove the `_ = regexp.MustCompile` line.

In `scrub`, return what the variable's `ReplaceAllString` makes of `s` with the replacement `$1.[redacted]`. There `$1` stands for what the pattern's first group matched, the customer ID, so each match keeps it and drops the secret.

If the test says `scrub` compiles its pattern, the `MustCompile` call is still inside `scrub`: move it into the variable. If a match comes out as `.[redacted]` with no customer, check the `$1` in the replacement.
