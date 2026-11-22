Look at the error before anything else. When it is not nil, the lines don't matter: one `fmt.Fprintf` to `stderr` with the `gctl: ` prefix and a newline, then return 1.

Only when the error is nil, loop over the lines and write each to `stdout` with `fmt.Fprintln`, then return 0. No lines at all is a normal, empty report.

If the test finds results on stdout after a failure, the lines are written before the error is checked: move the check to the top.
