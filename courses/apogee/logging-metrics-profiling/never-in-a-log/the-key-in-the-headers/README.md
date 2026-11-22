---
slug: the-key-in-the-headers
title: The Key in the Headers
type: code_tests
position: 0
difficulty: 7
topics: http, logging, security
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:38Z
---
Thirteen of the 14 bookingd lines Mirela found are DEBUG lines that dump each request's headers, `Authorization` and all. That header carries Furrowlight's API key, the secret its program sends to prove who it is, in plain text for anyone who reads logs in the ops console.

Redaction, replacing a secret with a placeholder such as `[redacted]`, keeps a line's shape and drops the secret. Redact a copy, since the request's own header must stay whole; `Header.Clone` copies every value:

```go
h := http.Header{}
h.Add("apogee-timestamp", "1833027183")
c := h.Clone()
c.Set("Apogee-Timestamp", "0")
fmt.Println(h.Get("Apogee-Timestamp"), c)
// 1833027183 map[Apogee-Timestamp:[0]]
```

## Your task

`safeHeader` returns `h` itself. Make it return a copy in which every value of `Authorization`, `Cookie`, `Set-Cookie` and `Apogee-Signature` reads `[redacted]`, and every other header stays as it is.

> [!NOTE]
> `Add`, `Set` and `Get` store and look up a name in its canonical form, so `apogee-timestamp` is kept as `Apogee-Timestamp`.
