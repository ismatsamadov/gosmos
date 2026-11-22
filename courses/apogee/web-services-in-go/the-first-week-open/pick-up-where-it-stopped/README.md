---
slug: pick-up-where-it-stopped
title: Pick Up Where It Stopped
type: code_tests
position: 2
difficulty: 9
topics: files, http, streams
created: 2026-11-13T06:35:18Z
modified: 2026-11-13T06:35:18Z
---
Thursday 3 February, 08:10. Hollis: "FURROW-A 0201 1402: FOUR TRIES, NEVER WHOLE." Tuesday's image is 10,912,500,000 bytes. Each try broke partway over Furrowlight's line, and each retry asked for the rest and got the whole file again, from its first byte.

A client that has part of a file asks for the rest with `Range: bytes=6000-`. A server that understands answers `206 Partial Content`, sends only those bytes, and says which in `Content-Range`. `http.ServeContent` does all of it for any `io.ReadSeeker`. It also takes the content's modification time: a client that resumes with `If-Range` and the date it saw gets the whole file again if it has changed since.

```go
motto := strings.NewReader("Every pass counts.")
http.ServeContent(w, r, "motto.txt",
	time.Time{}, motto)
// Range: bytes=6-  →  206 "pass counts."
```

## Your task

`sendFile` copies the whole file, whatever the request asks for. Answer with `http.ServeContent` instead, giving it the file's name and modification time from `file.Stat()`.
