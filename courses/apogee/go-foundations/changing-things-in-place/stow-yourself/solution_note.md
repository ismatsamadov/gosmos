With `*Station` receivers, `yst.Stow(...)` passes the address of `yst`, so the fields the methods set are the ones on `yst`, not on a copy that's thrown away.
