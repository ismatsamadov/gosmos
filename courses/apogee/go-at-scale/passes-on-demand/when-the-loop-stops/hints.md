`yield(p)` returns a bool: true while the loop wants more passes, false once it has stopped.

Test it with `if !yield(p)`. A `break` there would leave only the loop over one day's passes, and the outer loop would ask `predict` about the next day and call `yield` again.

Leave the whole function with `return`. The test breaks after 3 passes, and fails if `yield` is called again or if `predict` was asked about any day but 1 July.
