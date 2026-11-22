Make a copy with `slices.Clone` and sort the copy, never `passes`: sorting works in place, and the caller still needs the original order.

In the function you pass to `slices.SortFunc`, compare the two AOS times with their `Compare` method, which compares instants rather than clock readings. Give `cmp.Or` that result first and `strings.Compare` of the two IDs second.

Return the sorted copy. A tie on AOS that comes out in a different order on two runs means the IDs are not deciding it.
