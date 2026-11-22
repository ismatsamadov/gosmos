`Close` locks `h.mu` like every other method, sets `h.closed`, then ranges over `h.subs`, deleting each key and closing its channel, as `Publish` does with a full one. A second `Close` finds the map empty, so nothing is closed twice.

In `Subscribe`, check `h.closed` under the lock: when it is set, close the new channel and return it without storing it. The stream then sees a closed channel at once, says goodbye and returns.

In `newServer`, keep the server in a variable, call `RegisterOnShutdown` with the method value `b.hub.Close`, written without parentheses, and return the server.
