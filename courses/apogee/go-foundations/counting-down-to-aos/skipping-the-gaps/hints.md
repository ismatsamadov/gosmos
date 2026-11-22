Which minutes should count toward the average, and what do you divide by: the number of minutes, or the number of readings you kept? Think about a range where every minute is a gap.

Loop with a three-part `for` from `from` while the minute is under `to`. When `batteryAt` returns `lostFrame`, `continue`; otherwise add the reading to a sum and count it.

Shape: before the loop, a `float64` sum and an `int` count, both zero; after it, `if ___ == 0 {` returns 0, and otherwise the function returns the sum divided by `float64(___)`. For a `float64`, 0 divided by 0 is NaN, not 0.
