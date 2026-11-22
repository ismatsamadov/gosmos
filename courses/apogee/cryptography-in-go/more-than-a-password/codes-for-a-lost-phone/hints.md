In `Issue`, fill a second slice with `tokenDigest(codes[i])` as you make each code, keep that slice in `r.unused[user]`, and return the codes themselves for printing.

In `Use`, hash the typed code once with `tokenDigest`, then compare it with each kept digest using `subtle.ConstantTimeCompare`; the index of the match tells `slices.Delete` what to remove.

If a code from the sheet works twice, `Use` returns true without storing the shortened slice back in `r.unused[user]`.
