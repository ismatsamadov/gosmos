With a `*Station` parameter, `stow` receives the address of `yst` instead of a copy, so `s.Stowed = true` sets the field of the same `yst` that `main` prints.
