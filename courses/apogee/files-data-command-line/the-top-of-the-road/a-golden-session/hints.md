Open `testdata/help.golden` and `testdata/no-folder.golden`: each shows which stream gets which text, and the status to return. `[stdout]`, `[stderr]` and `[exit]` are the test's own markers, not output.

Both usages already end with a newline, so print the one `parseArgs` returns with `fmt.Fprint`, not `Fprintln`, or the transcript gains an empty line. Only a failure of the day itself starts with `gctl: `; a usage problem is printed as it is.

Check for `flag.ErrHelp` with `errors.Is` before the general check for any error, then hand the folder and format from the options to `planDay`. An error from `writeReport` is a failure too, with status 1.
