Lock `w.mu`, range over `w.conns` appending each connection to a slice, and unlock before writing anything.

`io.WriteString` writes the line and its newline to one watcher. When it returns an error, call `w.Remove` with that watcher: `Remove` takes the lock itself and closes the connection.

If `Broadcast` never returns while a watcher joins, the writes happen under the lock. If a watcher misses lines, a failed write ended the loop: remove that watcher and carry on with the next.
