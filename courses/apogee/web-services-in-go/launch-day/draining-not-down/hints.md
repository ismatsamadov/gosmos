The flag's `Store` method sets it, and `time.Sleep` keeps `run` waiting for `t.Drain` while the server, in its own goroutine, goes on answering.

In `health`, the flag's `Load` method reads it. Answer with `writeJSON` and `http.StatusServiceUnavailable`, and return before the line that answers 200.

Set the flag before the sleep, and sleep before `stop`: once `Shutdown` has begun, no new request can reach `/health` to see the flag.
