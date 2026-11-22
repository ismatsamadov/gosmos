---
slug: keep-the-line-open
title: Keep the Line Open
type: code_tests
position: 4
difficulty: 9
topics: http, loops, networking
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T19:06:41Z
---
Monday evening. Sunniva has timed the backhaul: each new connection costs a 0.7 s round trip, and the collector makes three requests per recording.

HTTP/1.1 keeps a connection open for the next request unless either side says `Connection: close`. The station server is stricter: it reads on only after a 2xx to a request that didn't say it, since a refused request may leave bytes on the line. Any other answer says `Connection: close`, and the server hangs up. `Peek` waits for a request's first byte without taking it:

```go
c.SetReadDeadline(time.Now().Add(time.Second))
b, err := br.Peek(1) // waits, keeps the byte
fmt.Printf("%q %v\n", b, err) // "G" <nil>
```

## Your task

`serveConn` in `server.go` answers one request, then closes. Make it loop: limit each `Peek` to `Idle` with the given `limit`, and answer with the given `send`, passing `last` true unless the answer is a 2xx and `closing(req)` is false. Stop when `send` reports the line can't go on.
