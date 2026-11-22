`SetWriteDeadline` works like `SetReadDeadline`: call it on `c.conn` with the instant `time.Now().Add(c.Timeout)`.

Set it at the top of `send`, before `writeFrame`, so every frame gets its own `c.Timeout` from the moment dishd starts writing it.

If the controller that never reads still holds `send` after a minute, the deadline set is a read deadline, or isn't on `c.conn`. If send 4 fails, the write deadline was set once, perhaps in `newClient`, and had passed by then.
