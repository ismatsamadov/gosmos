Each loop answers a question the `slices` package already answers: is it there, where is it, which is largest? The file needs to import `"slices"` first.

`slices.Contains(s, v)` reports whether `v` is in `s`, `slices.Index(s, v)` returns its position counted from 0, and `slices.Max(s)` returns the largest element.

Keep the three variables and their `Println` lines. Only the loops go: each variable gets its value from one call instead.
