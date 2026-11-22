---
slug: the-bytes-not-the-label
title: The Bytes, Not the Label
type: code_tests
position: 1
difficulty: 9
topics: bytes, http, security
created: 2026-11-13T06:35:18Z
modified: 2026-11-13T06:35:18Z
---
Wednesday, 15:00. Furrowlight's command file for Friday's pass over Salar Quieto is a PDF of the command list. Its part said `text/plain`, as everything Furrowlight's script sends does, and the service kept it.

A part's `Content-Type` says whatever its sender wrote. The station will send up the bytes, so the service looks at them: `http.DetectContentType` reads at most the first 512 and names what they look like.

```go
fmt.Println(http.DetectContentType(
	[]byte("<b>PING</b>")))
// text/html; charset=utf-8
```

## Your task

`checkCommands` accepts every file. Unless the sniffed type is exactly `text/plain; charset=utf-8`, return an `*APIError` with status 415, code `unsupported_media_type` and the message `command file must be plain text, not ` followed by the sniffed type.
