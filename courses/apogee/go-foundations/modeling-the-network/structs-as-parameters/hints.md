Salar Quieto's `LatDeg` is -23.4, but its line shows `23.4 S`. Which two things change when a latitude is negative?

An `if s.LatDeg < 0` can choose the letter and flip the sign. `fmt.Sprintf` with `%.1f` then builds the line from the fields.

Keep the letter and the latitude in two variables that start as `"N"` and `s.LatDeg`, change both inside the `if`, then pass them to `Sprintf` in the order the line shows them.
