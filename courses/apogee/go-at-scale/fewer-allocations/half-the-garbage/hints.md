Start with the scratch: take one from `getScratch` at the top, `defer` its `putScratch`, and append the passes to `*work` instead of to a slice that starts at nil. `assign` copies each pass into its contact, so nothing you return points into the scratch.

Make `seen` with `make` and `len(cands)` as its size, change its key type to `passKey`, and build each key from the satellite, the station and `p.AOS.Unix()`, as you did for the shadow run.

If the bytes are still over half, check that the copy lives in the scratch and that the scratch goes back: one that never does is made again on every call. If the contacts differ, sort the scratch and range over it for `assign`, never `cands` itself.
