Which call does exactly 50 km/h get, and exactly 72? A gale of 95 km/h is also above 50, so which limit must the first case test?

Write `switch {` with nothing after it, give each call a `case` holding a comparison on `windKmh`, and let `default` cover calm wind.

Shape: `case windKmh ___ ___:` returns `"stow"`, a second case returns the warning, and `default:` returns `"track"`.
