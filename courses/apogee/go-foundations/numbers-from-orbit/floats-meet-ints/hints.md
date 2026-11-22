In `float64(totalMB / passes)`, which happens first, the division or the conversion? Work out `1000 / 3` with two `int`s.

Wrap each `int` in `float64(...)` before the `/`, so the division itself happens between two `float64`s. Converting only one of them won't build.
