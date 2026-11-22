Start from a new, empty slice with room for every pass: `make` with a length of 0 and a capacity of `len(all)`. Append each pass whose `LOS` is after `now`. `After` says false when the two times are equal, so a pass ending this very second is gone, and one under way, past its AOS but not its LOS, stays.

Sort the new slice with `slices.SortFunc`. A `time.Time` has a `Compare` method that returns -1, 0 or +1, as `cmp.Compare` does for numbers and strings; when two AOS times compare as 0, the booking IDs decide, which is what `cmp.Or` says.

Keep at most 12 by slicing the sorted passes up to 12, or up to their length when fewer are left: the built-in `min` gives that end in one call. Add `"cmp"` to the imports of `board.go`.
