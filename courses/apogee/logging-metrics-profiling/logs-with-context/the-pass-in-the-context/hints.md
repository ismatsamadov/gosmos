`withPass` is one call: `context.WithValue` with the context it was given, the key `passKey{}` and the ID.

In `passFrom`, ask `ctx.Value(passKey{})` and turn the answer into a string with a type assertion, `.(string)`, using its two-result form so a missing value gives `""` instead of a panic.

If the test says an empty ID came back with `true`, `passFrom` returns the assertion's own `ok`: what it must report is whether the ID it found is not empty.
