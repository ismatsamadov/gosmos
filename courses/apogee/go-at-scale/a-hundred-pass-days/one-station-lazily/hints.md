`func(yield func(PassView) bool)` is the shape of an `iter.Seq[PassView]`: write the sequence as a function literal of that shape and pass it straight to `assign`. Inside, range over `stationPasses` for the station and window, with `passModel` as its model and `cfg.MinElevDeg` as the lowest elevation.

Keep a count of the candidates the loop has seen. When the count is a multiple of 1,000, starting at 0, return from the sequence if `ctx.Err()` isn't nil; skip a pass `stowedDuring` reports for `stows`, and return when `yield` says false.

After `assign`, check `ctx.Err()` once more: a sequence that stopped early leaves `assign` a short day, and those contacts must not come back as a plan. If the week's test says the model was asked for 4 July, the check isn't running inside the loop.
