Keep a counter that starts at 0 before the loop and goes up by one as the first thing inside it, so a note or blank line that reaches `continue` has already been counted.

Wrap the helper's error with `fmt.Errorf`, the counter and `%w`, in the form `line 7: ...`. With `%v` the text looks the same, but `errors.Is` can no longer find the `ParseFloat` error inside.
