The directive is a comment line of its own: `//go:embed` followed by the file's name. Put it right above a `var` declared at package level, outside any function.

Declare the variable as a `string`: `go build` fills it with the file's text. Then `stationTable` only has to return that variable.

If the test still gets an empty string, Go took the directive for an ordinary comment: write `//go:embed` with no space after the slashes, and check that `stationTable` returns the variable right below it.
