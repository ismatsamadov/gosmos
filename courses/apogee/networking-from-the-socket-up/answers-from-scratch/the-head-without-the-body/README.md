---
slug: the-head-without-the-body
title: The Head Without the Body
type: code_tests
position: 3
difficulty: 9
topics: http, networking
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T19:06:41Z
---
Monday afternoon. Before the collector fetches an 11 GB recording, Ingrid wants it to ask how big it is, and that answer must not bring the 11 GB with it.

A HEAD request asks for the answer a GET would get, minus its body: the same status and fields, `Content-Length` included, then nothing after the blank line. The client takes the length from the head and expects no body, so any byte that came would be read as the start of the next answer:

```go
res, err := client.Head(url)
fmt.Println(res.ContentLength, err)
// 4096 <nil>
```

## Your task

`route` now sends GET and HEAD of a recording to `recording`, in `recordings.go`, which answers HEAD with the body too. For HEAD, keep every field, `Content-Length` the full size, and leave the body nil, for the recording and for its 404.

> [!WARNING]
> Drop the body only after the fields are built: `Content-Length: 0` tells the collector the recording is empty.
