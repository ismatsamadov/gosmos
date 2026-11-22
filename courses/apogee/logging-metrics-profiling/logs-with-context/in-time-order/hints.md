`merge` already has the joined slice: keep it in a variable, sort it, and return it. `slices.SortStableFunc` takes the slice and a function of two `Line` values that returns an `int`.

The function compares the two lines' times, and `Time.Compare` already returns the -1, 0 or +1 that `SortStableFunc` wants from it.

If the test says the 01:48:30 lines are swapped, or passplan's four, the sort isn't stable: use `SortStableFunc`, not `SortFunc`, and compare times only, so equal times stay as they came.
