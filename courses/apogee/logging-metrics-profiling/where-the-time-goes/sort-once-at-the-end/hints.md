Start an empty `[]Planned` of the day's own, and add each group to it whole with `append` and `...`: no sorting yet.

After the loop, one call sorts the day: `slices.SortFunc` takes the slice and the comparing function, and `earlier` is already written, so it goes in by name, without calling it.

If the test says a group changed, the day began as one of the groups themselves, so sorting it moved their passes: start from an empty slice instead. If Brine-3 and Furrow-B come out in the wrong order, the sort compares AOS only; `earlier` breaks the tie by ID.
