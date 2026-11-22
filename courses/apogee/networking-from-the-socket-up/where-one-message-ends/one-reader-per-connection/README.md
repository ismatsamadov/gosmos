---
slug: one-reader-per-connection
title: One Reader per Connection
type: code_tests
position: 4
difficulty: 7
topics: networking, streams
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Thursday. To save a round trip on the backhaul, dishd sends HELLO and QUERY, its request for the dish's status, together, and the WELCOME and the STATUS come back glued in one read. In Wednesday's test, dishd lost the STATUS in which the simulated Dish 1 had stowed itself in the wind, and nobody saw the stow.

A `bufio.Reader` reads from its connection in big gulps, hands out what you ask for and keeps the rest in its buffer. Those bytes have left the connection: another reader on it will never see them. So a connection gets one reader for its whole life, and every read goes through it:

```go
br := bufio.NewReader(strings.NewReader(
	"first\nsecond\n"))
br.ReadString('\n')
fmt.Println(br.Buffered()) // 7: "second\n"
```

## Your task

`newClient` reads the banner through `c.br`, but `recv` makes a new `bufio.Reader` for every frame. Make `recv` read every frame through `c.br`.
