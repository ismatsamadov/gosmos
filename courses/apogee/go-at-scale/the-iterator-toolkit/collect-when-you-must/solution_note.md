`maps.Keys` yields the ten live codes one at a time, and `slices.Sorted` gathers them into a slice and sorts it, which needs every code at once. `strings.Join` then puts one space between them.
