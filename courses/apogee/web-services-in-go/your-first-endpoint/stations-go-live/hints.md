Start with the listener: `net.Listen("tcp", "127.0.0.1:0")` returns it and an error, and an error means there is nothing to serve, so return it.

Make an `http.Server` whose `Handler` field is `s.routes()`, and a channel of errors with room for one value. Start a goroutine that sends whatever `srv.Serve(ln)` returns into that channel, so `start` can return while the server keeps running.

If the test waits a whole second for `start`, `Serve` runs outside a goroutine; if `addr` isn't `127.0.0.1`, check the address you listen on. `ln.Addr().String()` is the address to return.
