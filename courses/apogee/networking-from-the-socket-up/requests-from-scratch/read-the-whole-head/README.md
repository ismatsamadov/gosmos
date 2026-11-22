---
slug: read-the-whole-head
title: Read the Whole Head
type: code_tests
position: 8
difficulty: 9
topics: errors, http, networking
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T19:06:41Z
---
Friday 26 May: the station server answers its first requests in the hut, each good one with the list of recordings and each bad one with its reason.

`readRequest` puts the pieces in order: the request line, where one over 1024 bytes is 414 `request line over 1024 bytes`; the header; exactly one Host, else 400 `missing host` or `more than one host`; a method of GET, HEAD or PUT, else 501 `method "PATCH" not implemented`; then the body's length. Counting a field's values tells one Host from two:

```go
h := Header{"host": {"127.0.0.1", "127.0.0.1"}}
fmt.Println(len(h["host"]), h.Get("host"))
// 2 127.0.0.1, 127.0.0.1
```

## Your task

`readRequest` in `request.go` returns an empty request. Check the head in that order, returning each error as it comes, then return the request with its method, target, header and length.

> [!NOTE]
> The test also sends heads to the given server on a loopback port; a head unfinished after 50 ms gets 408.
