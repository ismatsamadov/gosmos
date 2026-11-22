`maps.Keys(live)` is a sequence of the ten codes, in no fixed order. `slices.Sorted` takes a sequence like it and returns a sorted slice.

Put `slices.Sorted` of `maps.Keys(live)` where the list of five is, keep the `strings.Join` that prints it, and add `"maps"` and `"slices"` to the imports.

If the codes come out in a different order on each run, they aren't sorted yet: `slices.Collect` keeps the order `maps.Keys` yields them in, which changes from run to run. If you see five codes, `main` still prints the old list.
