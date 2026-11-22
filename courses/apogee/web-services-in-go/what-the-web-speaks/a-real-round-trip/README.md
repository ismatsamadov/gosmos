---
slug: a-real-round-trip
title: A Real Round Trip
type: code
position: 8
difficulty: 5
topics: defer, http, streams
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Friday. Before Cirrine's engineers use the stand-in themselves, Nkem wants to watch a request cross a real connection, the way theirs will.

`httptest.NewServer` serves a handler at `127.0.0.1`, the address by which your Mac reaches itself, behind a numbered door the system picks, called a port. `srv.URL` holds the whole address. `http.Get` sends a request there and returns the response, whose `Body` is a stream: read it, then close it.

```go
resp, err := http.Get(srv.URL + "/health")
if err != nil {
	return err
}
defer resp.Body.Close()
fmt.Println(resp.Status) // 404 Not Found
```

## Your task

`fetch` returns empty strings. Make it call `http.Get`, read the whole body with `io.ReadAll`, close it, and return `resp.Status` and the body, so `main` prints:

```text
200 OK {"code":"PEL","name":"Port Ellery"}
404 Not Found no station "ZZZ"
```

> [!NOTE]
> `http.Get` waits as long as the server takes. The stand-in answers in a millisecond or so; a slower service needs a client with a time limit.
