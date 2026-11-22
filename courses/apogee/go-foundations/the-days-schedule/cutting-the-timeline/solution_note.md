`min` keeps `n` within the slice's length, so neither bound runs past the end. `ids[:n]` takes the first `n` IDs, and starting at `len(readings)-n` with no end bound takes the last `n` readings.
