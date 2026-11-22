Count the spaces between `thistle` and `4182`. What does `strings.Split` hand back for the gap between two spaces that sit side by side?

The `strings` package has a function that splits around runs of white space and never makes empty pieces: `strings.Fields`.

`strings.Fields` takes only the string, with no separator, and returns a `[]string`, as `Split` does; nothing else in `main` changes.
