`SetReadDeadline` is a method of the connection, `c.conn`, and takes a `time.Time`: the instant after which reads fail. `time.Now().Add(c.Timeout)` is the instant `c.Timeout` from now.

Make the call where the `TODO` is, before `readFrame`, so every frame gets its own `c.Timeout` from the moment dishd starts waiting for it.

If exchange 4 fails, the deadline was set once, perhaps in `newClient`, and had passed by then. If the silent controller ends with `controller closed: idle` after 30 s, no deadline reached the read: set it on `c.conn`, inside `recv`.
