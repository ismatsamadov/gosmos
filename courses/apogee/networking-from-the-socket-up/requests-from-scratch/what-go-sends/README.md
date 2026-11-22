---
slug: what-go-sends
title: What Go Sends
type: code
position: 0
difficulty: 8
topics: http, networking, streams
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Monday 22 May. Before you write the station server, the server in the hut that will keep every recording until headquarters collects it, you watch what its clients send. The recorder, the hut's program that records each pass, is a Go program using `net/http`, and so is headquarters' collector, which fetches the recordings.

An HTTP request begins as text on a socket: a request line, a line per header field, then a blank line, each ending in CR LF, `\r\n`. `ReadString('\n')` on a `bufio.Reader` returns one line, CR LF included:

```go
br := bufio.NewReader(strings.NewReader(
	"YST-1 READY PROTO 1\r\n"))
line, _ := br.ReadString('\n')
fmt.Printf("%q\n", line)
// "YST-1 READY PROTO 1\r\n"
```

## Your task

`main` prints only the first line Go's client sends. Replace the `TODO` and the three statements under it with a loop that prints each line with `%q`, through the blank line, so it prints:

```text
"GET /recordings HTTP/1.1\r\n"
"Host: 127.0.0.1\r\n"
"User-Agent: Go-http-client/1.1\r\n"
"Accept-Encoding: gzip\r\n"
"\r\n"
```
