Loop from the second reading to the last, comparing each with the one before it, and keep a running total in a `uint64` that starts at 0. A loop that never runs, for fewer than two readings, returns that 0 by itself.

When a value is at least the one before, add the difference. When it is lower, the program restarted: add the new value itself, since it is what the program has counted from 0.

If the test shows 18446744073709551614, a subtraction went below 0: check that you subtract only when the new value is not lower than the old one.
