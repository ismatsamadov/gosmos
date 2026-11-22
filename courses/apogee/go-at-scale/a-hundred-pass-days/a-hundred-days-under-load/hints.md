Write the model inside `planStation`, where it can see `ctx`, as a function literal with `Predictor`'s parameters and result: it takes an `Elements`, a `Site`, `from` and `to`, and returns a `[]PassView`.

Its body checks `ctx.Err()` first and returns nil when the context has ended; otherwise it returns what `passModel` gives for the same arguments. Pass it to `stationPasses` where `passModel` was.

If a station still asks for 4 July, either `stationPasses` still gets `passModel` itself, or your function asks before it checks: the check comes first, for every day.
