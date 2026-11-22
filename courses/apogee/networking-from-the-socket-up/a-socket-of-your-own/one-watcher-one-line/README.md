---
slug: one-watcher-one-line
title: One Watcher, One Line
type: code_tests
position: 4
difficulty: 6
topics: defer, networking, streams
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T07:44:27Z
---
Wednesday. The first thing dishd's watch port will do is greet whoever connects with the dishes it knows, so a watcher can tell at once what it is watching.

`ln.Accept()` waits for a client and returns a `net.Conn` of its own for that client. What you write to it reaches the client, and closing it ends the stream: a client reading to the end then gets `io.EOF`. Closing the listener only stops new clients; a connection stays open until it is closed itself.

```go
data, err := io.ReadAll(conn)
// returns once the other end closes
fmt.Printf("%q %v\n", data, err)
// "dishd watching YST-1\n" <nil>
```

## Your task

`serveOne` returns `nil` at once. Make it accept one connection from `ln`, write the given `greeting(names)` to it and close it. When `Accept` fails, return its error as it is.

> [!TIP]
> `io.WriteString(w, s)` writes a string to any writer, a connection included.
