Put a `for` loop around the `Peek` and everything after it, and start each turn with `limit(c.SetReadDeadline, s.Limits.Idle)`; a zero `Idle` means no limit, which `limit` handles.

When `s.read` fails, answer with `s.send(c, failure(err), true)` and return. Otherwise route the request, and replace the starter's `writeResponse` with `send`: `a.status/100 == 2` is true for every 2xx, and `closing(req)` says whether the request asked to close.

`send` adds `Connection: close` to the last answer and returns false after it, so its result says whether the loop goes on. If the collector's three requests take three connections, `last` is always true; if the 404 test fails, the loop goes on after an error.
