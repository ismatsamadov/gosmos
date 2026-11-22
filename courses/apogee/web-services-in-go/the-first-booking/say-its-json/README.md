---
slug: say-its-json
title: Say It's JSON
type: code_tests
position: 5
difficulty: 6
topics: http
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Brinewatch's first test posted a booking as a form, `sat_id=brine-3&station=LAN`, and the service answered with a complaint about an invalid character `s`. Brinewatch's engineers had no idea what that meant.

`Content-Type` names a body's format, so a service can refuse a format it doesn't read before it reads a byte, with `415 Unsupported Media Type` and a message that says what to send instead. A missing header gives `mime.ParseMediaType` an empty string, which it refuses too:

```go
_, _, err := mime.ParseMediaType("")
fmt.Println(err) // mime: no media type
```

## Your task

`decodeJSON` decodes whatever arrives. Before it touches the body, parse the media type of `r`'s `Content-Type`, and unless it parses and is `application/json`, return a 415 `unsupported_media_type` `*APIError` saying:

```text
send the body as application/json
```

> [!TIP]
> It is the check that `isJSON` made in your first week, now on the service's side of the conversation.
