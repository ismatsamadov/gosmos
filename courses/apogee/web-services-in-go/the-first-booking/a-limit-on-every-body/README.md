---
slug: a-limit-on-every-body
title: A Limit on Every Body
type: code_tests
position: 4
difficulty: 7
topics: errors, http, streams
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
A Furrowlight test sent a booking with a 40 MB note by accident, and the service read every byte of it into memory before it answered.

`http.MaxBytesReader(w, r.Body, n)` wraps a body so that reading past `n` bytes fails with an `*http.MaxBytesError`, and the server closes the connection after its answer. The decoder hands that error on: check its type and answer `413 Request Entity Too Large`, as the body may be fine but too long.

```go
body := http.MaxBytesReader(w, r.Body, 1024)
_, err := io.ReadAll(body) // 2,000 bytes sent
fmt.Println(err)
// http: request body too large
```

## Your task

1. In `decodeJSON`, decode from `http.MaxBytesReader(w, r.Body, maxBodyBytes)` instead of `r.Body`.
2. In `bodyError`, which both `Decode` calls use, return a 413 `too_large` `*APIError` when `err` is an `*http.MaxBytesError`, saying:

```text
request body over 65536 bytes
```

> [!NOTE]
> A body of exactly 65,536 bytes is accepted: only a byte past the limit fails.
