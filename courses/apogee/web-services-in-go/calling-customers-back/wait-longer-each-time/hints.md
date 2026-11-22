Start from `2 * time.Second` and double it in a loop that runs `attempt - 1` times at most, and stops early once the wait has reached 5 minutes.

The built-in `min` caps it at 5 minutes. Turn the wait into a `float64` to multiply it by the factor, then back into a `time.Duration`.

Draw exactly one `r.Float64()` per call: a second draw shifts every later wait away from the test's. If attempt 70 gives 0, the wait was shifted past what a `time.Duration` holds; stop doubling at 5 minutes instead.
