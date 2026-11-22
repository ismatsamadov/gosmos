`append(buf, sat...)` appends a string's bytes to a byte slice, and `append(buf, '/')` appends one byte. Each returns the slice to keep using, so assign it back to `buf` every time.

For the padding, an orbit number below 10000 needs a `'0'` in front for each of 10000, 1000, 100 and 10 that is still larger than it: 8154 gets one, and 22379 none. Then `strconv.AppendInt` appends the number itself, as an `int64` in base 10.

If the test counts allocations, something still builds a string on the way: `fmt`, `strconv.Itoa` or `+` all make a new one. Every piece must go straight into `buf`.
