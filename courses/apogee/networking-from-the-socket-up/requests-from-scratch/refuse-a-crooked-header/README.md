---
slug: refuse-a-crooked-header
title: Refuse a Crooked Header
type: code_tests
position: 5
difficulty: 9
topics: errors, http, security
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T19:06:41Z
---
Mirela's second rule: when two programs could read a header line two ways, refuse it. To one program `Host : 127.0.0.1` names the host; to the next it is no Host at all.

A name is one or more token characters, letters, digits and marks such as `-` and `_`, as the given `isToken` checks, with nothing between it and the colon. A value may hold spaces and tabs but no other control character, such as a NUL byte, as the given `isFieldValue` checks: HTTP forbids them, and programs read them in different ways. A line with no colon, a crooked name or value, or one starting with a space or a tab, an old way to fold a field over two lines, is 400 `bad header line`:

```go
fmt.Println(isToken("Content-Type")) // true
fmt.Println(isToken("Host "))        // false
fmt.Println(isFieldValue("a\x00b"))  // false
```

## Your task

Make `readHeader` refuse those lines, a header line over 1024 bytes (431 `header line over 1024 bytes`) and a 33rd field (431 `more than 32 header fields`), each as a `RequestError`. Count lines, not names.

> [!WARNING]
> Don't trim the name: `Host :` would become `host`, and a crooked line would pass.
