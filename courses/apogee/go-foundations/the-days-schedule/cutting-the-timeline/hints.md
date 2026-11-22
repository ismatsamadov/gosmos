In a slice of length 3, the last 2 elements start at index 1. How do you get that 1 from the length and `n`?

The built-in `min(a, b)` returns the smaller of two numbers, which keeps `n` within `len`. Then slice with one bound left out.

`nextN` returns `ids[:___]`, and `lastN` returns `readings[___:]`.
