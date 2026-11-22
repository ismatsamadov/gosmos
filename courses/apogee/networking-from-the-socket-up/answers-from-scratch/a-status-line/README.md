---
slug: a-status-line
title: A Status Line
type: code
position: 0
difficulty: 8
topics: formatting, http, strings
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T07:44:27Z
---
Monday 19 June. Nkem hands you the station server's answers: a helper you were given wrote them until now, and this week you write them yourself, starting with the line every answer begins with.

That status line is the version, a space, the three-digit code, a space and its reason phrase, ended by CR LF. A client acts on the code; the reason is for people, and the given `reasons` table holds the words `net/http` uses:

```go
line := fmt.Sprintf("%s-%d READY\r\n", "YST", 2)
fmt.Printf("%q\n", line) // "YST-2 READY\r\n"
```

## Your task

`statusLine` in `response.go` returns the 200 line for every code. Build the line from `code` and its reason in `reasons`, so `main` prints:

```text
"HTTP/1.1 200 OK\r\n"
"HTTP/1.1 201 Created\r\n"
"HTTP/1.1 404 Not Found\r\n"
"HTTP/1.1 413 Request Entity Too Large\r\n"
```

> [!NOTE]
> Newer specifications call 413 "Content Too Large"; a client reads only the number, so either works.
