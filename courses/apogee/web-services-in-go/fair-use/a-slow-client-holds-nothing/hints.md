Each limit is a field of `http.Server` with a name ending in `Timeout`, and each takes a `time.Duration` straight from the matching field of `t`, a zero included.

`MaxHeaderBytes` is a plain `int` of bytes: 16 KiB is 16 × 1024.

If the connection that sent half a request stays open, `ReadHeaderTimeout` isn't set; if the idle one stays open, `IdleTimeout` isn't. If the 32 KiB header is answered 200, `MaxHeaderBytes` is still 0, which means Go's default of 1 MB.
