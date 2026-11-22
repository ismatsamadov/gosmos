Make the gzip reader from `r` and the tar reader from the gzip reader. Return the gzip error at once: a log that was never gzip'd fails right there.

Loop on `tr.Next()`. `io.EOF` means the parcel ended cleanly, so return the map; any other error goes back to the caller. Between two calls, reading from `tr` itself reads the current file.

Read each file whole with `io.ReadAll`, through `io.LimitReader` set to `limit+1`. More than `limit` bytes back means the file is too long; exactly `limit` is fine.
