---
slug: one-length-never-two
title: One Length, Never Two
type: code_tests
position: 5
difficulty: 9
topics: errors, http, security
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T19:06:41Z
---
Wednesday, Mirela's review: one upload with both `Content-Length` and `Transfer-Encoding: chunked`, which the station server reads by its length and the next program reads chunked, finding a second request hidden inside. "Your code accepts the request. Now show me the one it refuses."

So the station server decides a body's length in this order, the first failure answering. Both fields present is 400. The codings, the field's values joined with `, `, must end in `chunked`, or the body's end can't be found: 400. Any coding before that last one is one the server doesn't implement: 501. Chunked on GET or HEAD is 400. Else the body is chunked, length -1. Codings match in any letter case:

```go
strings.EqualFold("YST-1", "yst-1") // true
```

## Your task

`bodyLength` returns `contentLength`'s answer. Add those checks before it, with these messages:

```text
both content-length and transfer-encoding
bad transfer-encoding "gzip"
transfer-encoding "gzip, chunked" not implemented
GET has no body
```
