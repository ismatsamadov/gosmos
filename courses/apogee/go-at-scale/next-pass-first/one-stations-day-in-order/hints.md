Loop over `n.Elements` with its index: `n.Satellites[i]` is the same satellite, and `canServe(s, n.Satellites[i])` says whether the station shares one of its bands.

For each satellite it can serve, append `Filter(Passes(predict, e, s, w), keep)` to a slice of sequences, where `keep` reports whether a pass's `MaxElevDeg` is at least `minElevDeg`. Return `MergeBy(byAOS, seqs...)`.

If the test says a later day was predicted after an earlier one gave a candidate, passes are collected before the merge: give `MergeBy` the lazy sequences themselves.
