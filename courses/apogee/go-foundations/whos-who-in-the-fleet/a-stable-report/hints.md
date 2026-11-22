A `range` over `counts` hands you the IDs in any order. What could you range over instead, so the IDs arrive sorted?

Sort the IDs first, the way the example sorts the station codes, and range over that slice instead of the map. For each ID, build its line with `fmt.Sprintf` and `append` it to the slice you return.

The line for one ID is `fmt.Sprintf("%s: %d", ___, ___)`: the ID, then its count from the map.
