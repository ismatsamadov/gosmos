Find the bucket first: start from `len(h.Bounds)`, the bucket above every bound, and loop over `h.Bounds`, taking the first index whose bound is at least `d.Seconds()`.

Every field is atomic already, so each of the three changes is one call to `Add`: one for the bucket you found, one for the count, and the whole milliseconds of `d` for the sum.

If exactly 2 s lands in the "up to 5" bucket, the comparison leaves a bound out: a wait equal to a bound belongs in that bound's bucket, so compare with `<=`.
