`ln.Accept()` returns the new connection and an error. Return the error when there is one: on a closed listener there is no connection to use.

Once you have the connection, defer its `Close` and write `greeting(names)` to it with `io.WriteString`, returning the write's error.

If the watcher read the greeting and then an i/o timeout, the connection never closed: closing `ln` isn't enough, because each connection is closed on its own.
