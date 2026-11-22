The lookup gives you two values, the `command` and a `bool`. When the bool is false, keep the unknown-command line; when it is true, the command's `run` field is a function you call like any other, and `words[1:]` holds the words after the name.

`cmdHelp` takes the same parameters as a command's `run` and returns an `error` too. A map's order changes from run to run, so it must sort before it prints. Get the names in order with `slices.Sorted` over `maps.Keys` of the table, or append the keys to a slice and call `sort.Strings`, then look each name up again for its help text.

Each line of `help` is the name through `%-8s`, two spaces, then the help text and a newline. In the table, the entry for `"help"` holds `cmdHelp` itself as its `run` value, without parentheses. Add `"fmt"`, `"maps"` and `"slices"` (or `"sort"`) to the imports of `commands.go`.
