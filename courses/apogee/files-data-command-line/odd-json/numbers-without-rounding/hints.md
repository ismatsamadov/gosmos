Call `UseNumber` on the decoder before `Decode`. Afterwards every number in the map is a `json.Number` instead of a `float64`.

Read the `id` entry with a two-result type assertion to `json.Number`. A missing entry or an ID written as text fails the assertion, and that is your error case.

Return `Int64`'s two results as they are: its error already says whether the digits held a fraction or didn't fit in an `int64`.
