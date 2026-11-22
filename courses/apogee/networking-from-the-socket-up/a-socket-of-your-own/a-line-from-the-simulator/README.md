---
slug: a-line-from-the-simulator
title: A Line from the Simulator
type: code
position: 0
difficulty: 6
topics: networking, streams
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T07:44:27Z
---
Monday 6 March, in Ystmark's one building, the hut. The computers that will turn its dishes, the new dish controllers, wait in their crate; the maker's simulator, a program that behaves like one, runs on the bench. Sunniva: "Minus 18. What does it say?"

`net.Dial("tcp", addr)` connects to whatever listens at `addr`, returning a `net.Conn`. Its `Read` waits for bytes, copies them into your slice and returns how many, `n`; only `buf[:n]` holds them:

```go
r := strings.NewReader("Wind 64, gusts 81")
buf := make([]byte, 8)
n, err := r.Read(buf)
fmt.Printf("%d %q %v\n", n, buf[:n], err)
// 8 "Wind 64," <nil>
```

## Your task

`main` dials the simulator but never reads. Replace the `TODO` with `conn.Read(buf)`, setting `n` and printing any error, so it prints:

```text
21 bytes: "YST-1 READY PROTO 1\r\n"
```

> [!NOTE]
> `%q` shows the bytes you can't see: a carriage return, `\r`, and a newline, `\n`. Page one of the manual mentions only the `\n`.
