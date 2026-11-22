Split the line with `strings.Fields`. A good line has at least five fields: `CMD`, the satellite, the number, the operation and `sig=`. Everything between the operation and the last field is an argument, and none of them may start with `sig=`.

`strings.CutPrefix` on the last field gives the signature and tells you whether it started with `sig=`; check that exactly 64 lowercase hex digits follow. `strconv.ParseUint(f[2], 10, 64)` reads the number. Any failure returns `bad command line` before the guard is touched.

Don't check the signature or the number yourself: `g.Accept` builds the canonical form, verifies, compares and records, in that order. Return its error as it is, so that `errors.Is` still finds `ErrForged` or `ErrReplay`.
