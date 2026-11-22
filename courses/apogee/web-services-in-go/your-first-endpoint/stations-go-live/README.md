---
slug: stations-go-live
title: Stations Go Live
type: code_tests
position: 7
difficulty: 6
topics: goroutines, http
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Friday 19 November. Yusra promised Cirrine an address today for the booking service's first test run, which opens the beta: the weeks in which customers try the service before launch.

An `http.Server` serves a handler on a listener. `net.Listen("tcp", "127.0.0.1:0")` opens one on this Mac only, on a port the system picks, and `ln.Addr()` tells you which. `srv.Serve(ln)` runs until the server closes, so it goes in a goroutine.

```go
ln, err := net.Listen("tcp", "127.0.0.1:0")
if err != nil {
	return err
}
fmt.Println(ln.Addr()) // 127.0.0.1:52134
```

## Your task

`start` returns an error. Listen on `127.0.0.1:0`, make an `http.Server` whose `Handler` is `s.routes()`, run `srv.Serve(ln)` in a goroutine that sends its result on a channel, `served`, and return the listener's address and `stopFunc(srv, served)`.

> [!WARNING]
> `":0"` would listen on every network this Mac is on, and macOS may ask whether to allow that; `127.0.0.1` stays on the Mac.
