Which names come out too long, and what do they have that `Ystmark` doesn't?

The `unicode/utf8` package counts runes; add it to the imports next to `fmt`.

`utf8.RuneCountInString` takes a string and returns how many runes it holds, as an `int`, the way `len` returns its bytes.
