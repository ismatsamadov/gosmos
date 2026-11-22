Each `return` inside the checks ends `validate` at the first problem. What could you do with each problem instead?

Append each problem to a `[]error`, then return `errors.Join` with that slice followed by `...`.

No special case is needed for a healthy frame: `errors.Join` of an empty slice is `nil`.
