Read `Accept` against the order in the readme: it parses, compares the number and records it at once, then opens with whatever `s.Keys` holds under the line's key ID: nothing, or a key the satellite may no longer use.

After `parseSealed`, refuse a line whose `SatID` isn't `s.ID`, wrapping `ErrWrongSat` with `%w` and the satellite with `%q`. Then take the keys from `s.keyFor(sealed.KeyID)`, which returns the right error by itself, and only then compare the number with `s.Last`.

Call `openCommand` with `keys.Seal` and keep its results, so that `s.Last = sealed.Seq` comes after its error check, right before the return. If `TestForgedMaximum` still fails, something sets `s.Last` before `openCommand` has succeeded.
