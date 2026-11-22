---
slug: only-on-this-mac
title: Only on This Mac
type: code_tests
position: 3
difficulty: 6
topics: networking
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T07:44:27Z
---
The dish service will have a watch port: a port anyone can connect to and watch both dishes, a line per dish every second. Sunniva wants to try one on the bench computer without opening anything to the station's network.

`net.Listen("tcp", "127.0.0.1:0")` listens on `127.0.0.1`, the address by which a computer reaches itself, so nothing outside your Mac can connect. Port `0` asks the system for any free port, and the listener's `Addr` says which one it picked:

```go
host, port, err := net.SplitHostPort(
	ln.Addr().String())
if err != nil {
	return err
}
fmt.Println(host, port) // 127.0.0.1 52113
```

## Your task

`listenLocal` listens on `127.0.0.1:7300`, a fixed port that only one listener at a time can have. Let the system pick the port instead.

> [!WARNING]
> Keep the host. `:0` alone listens on every network the Mac is on, the station's too, and macOS may ask whether to allow it.
