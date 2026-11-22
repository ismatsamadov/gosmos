The receiver `(p Pass)` makes `Duration` a method of `Pass`, so `p.Duration()` compiles and returns `p.LOS.Sub(p.AOS)`, which `%v` prints as `9m12s` for Pipit-1's pass.
