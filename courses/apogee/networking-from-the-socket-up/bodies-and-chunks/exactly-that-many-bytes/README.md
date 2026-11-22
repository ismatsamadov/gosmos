---
slug: exactly-that-many-bytes
title: Exactly That Many Bytes
type: code_tests
position: 0
difficulty: 8
topics: errors, http, streams
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T19:06:41Z
---
Monday 5 June. The recorder, the program in the hut's rack that records each pass, uploads its first test pass to the station server: 4,096 bytes of Pipit-2, zeros, CRs and LFs, all to be kept.

A body with `Content-Length` is exactly that many bytes, read through the `bufio.Reader` that read the head, whose buffer may hold the first of them. `io.ReadFull` fills a slice however the reads arrive, or returns `io.ErrUnexpectedEOF` if the stream ends partway:

```go
buf := make([]byte, 8)
n, err := io.ReadFull(
	strings.NewReader("Wind 64"), buf)
fmt.Println(n, err) // 7 unexpected EOF
```

## Your task

`readBody` in `body.go` reads lines until it has `n` bytes. Read exactly `n` with `io.ReadFull` instead, and when the stream ends early, return a `*RequestError` with status 400 and the count:

```text
body ended after 812 of 4096 bytes
```

> [!NOTE]
> From today the given server reads the body into `req.Body` before routing, its length decided by `bodyLength`.
